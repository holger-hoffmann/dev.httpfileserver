When microservices need to exchange data, they can do so via files with well defined content.

In the cloud, this is done via storage buckets.

As we do not want each microservice to have access to the storage bucket and we do not want the microservice to know how to access the storage bucket, we can work with signed URLs.

The microservice just has to know how to download files via HTTP GET and upload files via HTTP PUT.

Tho emulate this behaviour in local development, this fileserver serves files via HTTP GET and you can upload a file via HTTP PUT.