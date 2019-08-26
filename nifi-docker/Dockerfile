FROM apache/nifi:1.9.2
USER root
ADD start.sh ${NIFI_BASE_DIR}/scripts/start.sh
RUN chown nifi:nifi ${NIFI_BASE_DIR}/scripts/start.sh
USER nifi
