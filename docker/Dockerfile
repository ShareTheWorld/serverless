# scheduler server image
FROM golang:1.13.7

ARG port=10600
#ARG server_dir=/aliyuncnpc/scheduler
#ARG config_dir=${server_dir}/config
#ARG log_dir=${server_dir}/log

#ENV SERVICE_PORT=${port}
#
#RUN mkdir -p ${server_dir}
#RUN mkdir -p ${config_dir}
#RUN mkdir -p ${log_dir}

COPY ./main /

EXPOSE ${port}

ENTRYPOINT ["/main"]
