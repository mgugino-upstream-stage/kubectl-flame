FROM golang:1.14-buster as agentbuild
WORKDIR /go/src/github.com/VerizonMedia/kubectl-flame
ADD . /go/src/github.com/VerizonMedia/kubectl-flame
RUN go get -d -v ./...
RUN cd agent && go build -o /go/bin/agent

FROM redhat/ubi8
WORKDIR /app
RUN yum install -y dotnet-sdk-6.0 jq && yum clean all
RUN curl -L -o dotnet-trace https://aka.ms/dotnet-trace/linux-x64 && chmod +x dotnet-trace
RUN curl -L -o containerd.tar.gz https://github.com/containerd/containerd/releases/download/v1.5.7/containerd-1.5.7-linux-amd64.tar.gz \
  && tar xf containerd.tar.gz && rm containerd.tar.gz
COPY --from=agentbuild /go/bin/agent /app/agent
COPY /hack/containertmp.sh /app/containertmp.sh
RUN chmod +x /app/containertmp.sh
CMD [ "/app/agent" ]
