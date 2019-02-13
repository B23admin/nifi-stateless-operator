FROM apache/nifi:1.8.0
USER root
ADD start.sh ${NIFI_BASE_DIR}/scripts/start.sh
RUN chown nifi:nifi ${NIFI_BASE_DIR}/scripts/start.sh
USER nifi
