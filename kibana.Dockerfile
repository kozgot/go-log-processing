FROM docker.elastic.co/kibana/kibana:7.9.3

RUN /usr/share/kibana/bin/kibana-plugin install https://github.com/walterra/kibana-milestones-vis/releases/download/v7.9.3/kibana_milestones_vis-7.9.3.zip --allow-root
