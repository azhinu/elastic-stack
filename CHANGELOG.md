# 6 Jun 2022
- **Kibana:**
  - Updated monitoring settings to use `log-internal-elastic` index in monitoring tab
- **Filebeat:**
  - JSON decoding processor now will be used only if there is no event.dataset. Used to fix Filebeat modules processing.

# 23 May 2022
- Changed hostfs path
- Elastic 8.2
- **Beats:**
  - Some fields now will be dropped
  - Added JSON message processing
  - Dashboards setup now disabled

# 25 Apr 2022

- **Beats:**
  - Metricbeat: Added environment variables to respect container hostfs option.
  - Filebeat: Changed HTTP metrics server hostname.
# 23 Apr 2022

- Updated to 8.1.3
- Docker compose now using extension fields and override extension.
- Updated `README.md`
- **Elastic:**
  - Added environment variables to get ability to run in custer mode.
  - Added labels to filebeat logging.
- **Kibana:**
  - Added custom elastic package registry to use Fleet in isolated environment.
  - Updated healthcheck params.
- **Logstash:**
  - Commented out.
- **Elastic Agent:**
  - Added Elastic Agent and Fleet services.
- **Beats:**
  - Added Filebeat and Metricbeat for cluster monitoring. Internal monitoring is depricated.
- **Elastic package registry:**
  - Added local package registry.

# 19 Mar 2022

- Updated README
- Updated logstash user
- Updated healthcheck


# 18 Mar 2022

- Updated to 8.1.0.
- Initial setup become easier a bit and now can be fully automated.
- Now CA using PEM cert instead of PKCS#12 to get rid of useless certificate doubles.
  Config files still are available.
- Now Elasticsearch using only one PKCS#12 certificate for HTTP and TCP.
- Now passwords are stored in .env file.
- Now service configs are stored as environment variables.
- Fixed healthcheck args.
- Some small fixes.

## 24 Dec 2021
- Init.
