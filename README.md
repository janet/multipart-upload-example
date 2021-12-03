# multipart-upload-example

aws-sdk-go-v2 multipart upload poc

There are 2 examples of uploading multipart with aws-sdk-go-v2:

1. Manually, starting by calling s3Client.CreateMultipartUpload
2. S3 Manager, a feature implemented by S3

## How to run the code

1. git clone the repo
2. chdir into either the manual_multipart or manager directories
3. find a file that you would like to upload and update the const variables at the top
4. within the example directory, run:

        go run .
5. look in the s3 bucket and find the uploaded file