package job

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/VerizonMedia/kubectl-flame/cli/cmd/data"
	"github.com/VerizonMedia/kubectl-flame/cli/cmd/version"
)

type dotnetCreator struct{}

func (p *dotnetCreator) create(targetPod *apiv1.Pod, cfg *data.FlameConfig) (string, *batchv1.Job, error) {
	id := string(uuid.NewUUID())
	var imageName string
	var imagePullSecret []apiv1.LocalObjectReference
	args := []string{
		id,
		string(targetPod.UID),
		cfg.TargetConfig.ContainerName,
		cfg.TargetConfig.ContainerId,
		cfg.TargetConfig.Duration.String(),
		string(cfg.TargetConfig.Language),
		string(cfg.TargetConfig.Event),
	}

	if cfg.TargetConfig.Pgrep != "" {
		args = append(args, cfg.TargetConfig.Pgrep)
	}

	if cfg.TargetConfig.Image != "" {
		imageName = cfg.TargetConfig.Image
	} else {
		imageName = fmt.Sprintf("%s:%s-dotnet", baseImageName, version.GetCurrent())
	}

	if cfg.TargetConfig.ImagePullSecret != "" {
		imagePullSecret = []apiv1.LocalObjectReference{{Name: cfg.TargetConfig.ImagePullSecret}}
	}

	commonMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("kubectl-flame-%s", id),
		Namespace: cfg.TargetConfig.Namespace,
		Labels: map[string]string{
			"kubectl-flame/id": id,
		},
		Annotations: map[string]string{
			"sidecar.istio.io/inject": "false",
		},
	}

	resources, err := cfg.JobConfig.ToResourceRequirements()
	if err != nil {
		return "", nil, fmt.Errorf("unable to generate resource requirements: %w", err)
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: commonMeta,
		Spec: batchv1.JobSpec{
			Parallelism:             int32Ptr(1),
			Completions:             int32Ptr(1),
			TTLSecondsAfterFinished: int32Ptr(5),
			BackoffLimit:            int32Ptr(2),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: commonMeta,
				Spec: apiv1.PodSpec{
					HostPID: true,
					Volumes: []apiv1.Volume{
						{
							Name: "host-filesystem",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
					},
					ImagePullSecrets: imagePullSecret,
					InitContainers:   nil,
					Containers: []apiv1.Container{
						{
							ImagePullPolicy: apiv1.PullAlways,
							Name:            ContainerName,
							Image:           imageName,
							Command:         []string{"/app/agent"},
							Args:            args,
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "host-filesystem",
									MountPath: "/host",
								},
							},
							SecurityContext: &apiv1.SecurityContext{
								Privileged: boolPtr(true),
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{"SYS_PTRACE"},
								},
							},
							Resources: resources,
						},
					},
					RestartPolicy: "Never",
					NodeName:      targetPod.Spec.NodeName,
				},
			},
		},
	}

	return id, job, nil
}
