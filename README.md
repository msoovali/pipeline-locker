# Pipeline status checker and locker
If there is a need to hold back pipeline from getting deployed over then this project comes to help.
## How it works?
Pipeline-Locker helps to hold locked/unlocked state of pipelines. Pipeline consists of 2 identifiers: project name and pipeline deploy environment. Until bash and curl is available to make requests inside the pipeline and pipeline fails on exit code 1 then it doesn't matter which CI runner is in use.
## How to set it up?
1. Get pipeline-locker up and running using GO or Docker image.
2. Add bash script with curl request to check status of pipeline into CI. Take a look at [pipeline-lock-checker.sh](https://www.github.com/msoovali/pipeline-locker/blob/master/pipeline-lock-checker.sh) file. No need to worry anymore if someone accidentally tries to deploy Your deployment over, their pipeline fails if environment is locked.
3. Profit
## Pipeline-Locker roadmap
1. ~~Implement redis support aside to application memory storage, so it is possible to have more than 1 replica and state remains on application restart. Make it configurable.~~ ✅
2. Add config to predefine pipelines and option to select pipelines from dropdown list.
3. Add possibility to call lock API from CI safely with authorization and make secret configurable.
4. Improve UI

## Configurable environment variables
|Key                        |Default       |Description                                                                                             |
|---------------------------|--------------|--------------------------------------------------------------------------------------------------------|
|ADDR                       |:8080         |Service ip:port                                                                                         |
|ALLOW_OVERLOCKING          |false         |Allow to lock already locked pipeline                                                                   |
|PIPELINES_CASE_SENSITIVE   |true          |Project and environment case sensitivity                                                                |
|REDIS_VERSION              |0             |Redis version. Default 0 means disabled and in-memory data store is used. Supported redis versions: 6, 7|
|REDIS_ADDR                 |localhost:6379|Redis ip:port                                                                                           |
|REDIS_USERNAME             |              |Redis username                                                                                          |
|REDIS_PASSWORD             |              |Redis password                                                                                          |