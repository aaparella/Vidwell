Vidwell
===

Vidwell is a video sharing platform, much like YouTube. Built entirely as an academic exercise, to see how a Go web application could be structured cleanly and efficiently. It is also an experiment to see how much of the infrastructure of the application can be made with Go. The object storage used for pictures and videos is done through Minio (or any S3 compatible object storage service), and ideally the database would be in Go as well though this has proved to be more challenging. 

Planned Features
===

* User accounts and subscriptions
* Email notifications
* Comments
* Video tags and search

How to Run
===

VidWell can easily be run locally. Simply modify the config.ini file with the URL of the object storage service, and the respective keys. Also set the database url appropriately. Then the application can simply be built and run as a normal go application. The following is an example config.ini file, with the values annotated.

```yaml
[Storage]

Endpoint=        # Object store endpoint 
AccessKeyID=     # Object store credentials
SecretAccessKey= # Object store credentials
UseSSL=          # Connect to object storage using SSL
Database=        # URL passed to gorm.Open()
DatabaseLog=     # Control log level of database results

[Rendering]
TempaltesDir=./views # Directory containing templates

[Session]
Key=             # Key used to encrypt cookies
```
