FROM docker.elastic.co/kibana/kibana:6.4.3
RUN rm /usr/share/kibana/config/kibana.yml
COPY kibana.yml /usr/share/kibana/config
CMD ["/usr/local/bin/kibana-docker"]