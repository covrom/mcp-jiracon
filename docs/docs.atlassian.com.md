# [Atlassian Developers](https://developer.atlassian.com/)

# Confluence REST API Documentation 

This document describes the REST API and resources provided by Confluence. The REST APIs are for developers who want to integrate Confluence into their application and for administrators who want to script interactions with the Confluence server. 

Confluence's REST APIs provide access to resources (data entities) via URI paths. To use a REST API, your application will make an HTTP request and parse the response. The response format is JSON. Your methods will be the standard HTTP methods like GET, PUT, POST and DELETE. 

Because the REST API is based on open standards, you can use any web development language to access the API.

## Structure of the REST URIs 

URIs for the Confluence REST API resource have the following structure:

With context: http://host:port/context/rest/api/resource-name

Or without context: http://host:port/rest/api/resource-name

Example with context: [http://example.com:8080/confluence/rest/api/space/ds](http://example.com:8080/confluence/rest/api/space/ds)

Example without context: [http://confluence.myhost.com:8095/rest/api/space/ds](http://confluence.myhost.com:8095/rest/api/space/ds)

## How to use expansion in the REST APIs 

In order to minimise network traffic from the client perspective, our API uses a technique called expansion.

You can use the `expand` query parameter to specify a comma-separated list of entities that you want expanded, identifying each entity by a given identifier. For example, the value `space,attachments` requests the expansion of entities for which the expand identifier is `space` and `attachments`. 

You can use the `.` dot notation to specify expansion of entities within another entity. For example `body.view` would expand the content `body` and expand the `view` rendering of it. 

## Resources

### [api/accessmode](#api/accessmode)

Resource for plugins to check the access mode set on Confluence.

#### [Get access mode status](#api/accessmode-getAccessModeStatus)`GET /rest/api/accessmode`

Returns the access mode status for Confluence.

- `http://example.com/confluence/rest/api/accessmode` 

##### Responses

  - *application/json*

### [api/admin/group](#api/admin/group)

REST API for \{@link Group\} management actions that require elevated privileges.

#### [Create](#api/admin/group-create)`POST /rest/api/admin/group`

Create the given \{@link Group\} identified by name.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/group#","title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}
```

##### Responses

  - *application/json*

#### [Delete](#api/admin/group-delete)`DELETE /rest/api/admin/group/{groupName}`

Delete the given \{@link Group\} identified by groupName.

##### Responses

  - *application/json*

### [api/admin/user](#api/admin/user)

REST API for \{@link User\} management actions that require elevated privileges.

#### [Create user](#api/admin/user-createUser)`POST /rest/api/admin/user`

Creates a new active user. 

One of the following options could be used 

1. Create a user with a specified password. The userName, fullName, email and password needs to be specified 

2. Create a user with an email notification to the user. The userName, fullName, email and notifyViaEmail (true) needs to be specified 

The response would be a JSON returning the User Key for the newly created user 

**Requirements** 

- The userName should not be null or blank 

- The userName should not contain any of these characters \\ , + \< \> ' " 

- The userName should not contain any whitespace characters 

- The userName should not be "anonymous" 

- The userName should not contain any upper case characters 

- The fullName should not be null or blank 

- The fullName should not contain any of these characters \< \> 

- The fullName should not be "anonymous" 

- The email should not be null or blank 

- The email should be a valid email address 

- If notifyViaEmail is false then the password should not be null or blank 

- If notifyViaEmail is true then the password should not be specified 


##### Request

###### Example

```
{"userName": "user1", "fullName": "someone", "email" : "test1@test.com", "password" : "password", "notifyViaEmail": "false"}


 {"userName": "user1", "fullName": "someone", "email" : "test1@test.com", "notifyViaEmail": "true"}
```

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/user-details-for-creation#","title":"User Details For Creation","type":"object","properties":{"userName":{"type":"string"},"fullName":{"type":"string"},"email":{"type":"string"},"password":{"type":"string"},"notifyViaEmail":{"type":"boolean"}},"additionalProperties":false}
```

##### Responses

  - **Status 201** - *application/json*

###### Example

  - ```
{"userKey": "4028d6768851c8900188556ba6000014"}
```

  - Returns a full JSON representation of the UserKey 

  - **Status 400**

  - returned if any error occurs while creating the user

  - **Status 401**

  - returned if an anonymous (or unauthenticated) user tries to create a user

  - **Status 402**

  - returned if no more licenses available to create a user

  - **Status 403**

  - returned if user does not have enough permission to create a user

  - **Status 409**

  - returned if the user with the same userName already exists

#### [Delete](#api/admin/user-delete)`DELETE /rest/api/admin/user/{username}`

Delete the given \{@link User\} identified by username. This action is processed asynchronously. 

##### Responses

  - *application/json*

#### [Disable](#api/admin/user-disable)`PUT /rest/api/admin/user/{username}/disable`

Disable the given \{@link User\} identified by username. This method is idempotent i.e. if the user is already disabled then no action will be taken. 

##### Responses

  - *application/json*

#### [Enable](#api/admin/user-enable)`PUT /rest/api/admin/user/{username}/enable`

Enable the given \{@link User\} identified by username. This method is idempotent i.e. if the user is already enabled then no action will be taken. 

##### Responses

  - *application/json*

#### [Change password](#api/admin/user-changePassword)`POST /rest/api/admin/user/{username}/password`

Change the password for the user identified by the username 

Validation rules : The new password should not be null or blank 

##### Request

###### Example

```
http://yourhost:8080/confluence/rest/api/admin/user/someuser/password
```

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/credentials#","title":"Credentials","type":"object","properties":{"password":{"type":"string"}},"additionalProperties":false}
```

##### Responses

  - **Status 204**

  - returned if the password changed successfully

  - **Status 400**

  - returned if any error occurs while changing user password. Refer to the validation rules above.

  - **Status 403**

  - returned if user does not have enough permission to change another user's password. User should be a System admin

  - **Status 404**

  - returned if user with specified userName not found

### [api/audit](#api/audit)

#### [Get audit records](#api/audit-getAuditRecords)`GET /rest/api/audit`

Fetch a paginated list of AuditRecord instances dating back to a certain time

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `startDate` | *string* |  |
| `endDate` | *string* |  |
| `start` | *int* | where to start within results set |
| `limit` | *int* Default: `1000` | the maximum results to fetch |
| `searchString` | *string* |  |

##### Responses

  - *application/json*

#### [Store record](#api/audit-storeRecord)`POST /rest/api/audit`

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/audit-record#","title":"Audit Record","type":"object","properties":{"author":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"remoteAddress":{"type":"string"},"creationDate":{"type":"string"},"summary":{"type":"string"},"description":{"type":"string"},"category":{"type":"string"},"sysAdmin":{"type":"boolean"},"affectedObject":{"$ref":"#/definitions/affected-object"},"changedValues":{"type":"array","items":{"title":"Changed Value","type":"object","properties":{"name":{"type":"string"},"oldValue":{"type":"string"},"newValue":{"type":"string"}},"additionalProperties":false}},"associatedObjects":{"type":"array","items":{"$ref":"#/definitions/affected-object"}}},"definitions":{"affected-object":{"title":"Affected Object","type":"object","properties":{"name":{"type":"string"},"objectType":{"type":"string"}},"additionalProperties":false}},"additionalProperties":false,"required":["sysAdmin"]}
```

##### Responses

  - *application/json*

#### [Export](#api/audit-export)`GET /rest/api/audit/export`

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `startDate` | *string* |  |
| `endDate` | *string* |  |
| `searchString` | *string* |  |
| `format` | *string* Default: `csv` |  |

##### Responses

  - *application/zip*

  - *text/csv*

#### [Get retention period](#api/audit-getRetentionPeriod)`GET /rest/api/audit/retention`

Fetches the current retention period

##### Responses

  - *application/json*

#### [Set retention period](#api/audit-setRetentionPeriod)`PUT /rest/api/audit/retention`

Set the retention period to a new value. Can throw ServiceException if the retention period is too long

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/retention-period#","title":"Retention Period","type":"object","properties":{"units":{"type":"string","enum":["Nanos","Micros","Millis","Seconds","Minutes","Hours","HalfDays","Days","Weeks","Months","Years","Decades","Centuries","Millennia","Eras","Forever"]},"number":{"type":"integer"}},"additionalProperties":false,"required":["number"]}
```

##### Responses

  - *application/json*

#### [Get audit records](#api/audit-getAuditRecords)`GET /rest/api/audit/since`

Fetch a paginated list of AuditRecord instances dating back to a certain time

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `number` | *long* Default: `3` | the amount of time periods |
| `units` | *string* | the units to use for the time periods eg. 'days', 'months' etc |
| `start` | *int* | where to start within results set |
| `limit` | *int* Default: `1000` | the maximum results to fetch |
| `searchString` | *string* |  |

##### Responses

  - *application/json*

### [api/backup-restore](#api/backup-restore)

#### [Create site backup job](#api/backup-restore-createSiteBackupJob)`POST /rest/api/backup-restore/backup/site`

Starts the new site backup job.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/site-backup-settings#","title":"Site Backup Settings","type":"object","properties":{"skipAttachments":{"type":"boolean"},"keepPermanently":{"type":"boolean"},"fileNamePrefix":{"type":"string"}},"additionalProperties":false,"required":["skipAttachments","keepPermanently"]}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the site backup job 

  - **Status 400**

  - Returned if invalid settings provided

  - **Status 403**

  - Returned if user doesn't have permission to create site backups

#### [Create space backup job](#api/backup-restore-createSpaceBackupJob)`POST /rest/api/backup-restore/backup/space`

Creates new space backup job and adds it to the queue.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/space-backup-settings#","title":"Space Backup Settings","type":"object","properties":{"spaceKeys":{"type":"array","items":{"type":"string"}},"keepPermanently":{"type":"boolean"},"fileNamePrefix":{"type":"string"}},"additionalProperties":false,"required":["keepPermanently"]}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the space backup job 

  - **Status 400**

  - Returned if invalid settings provided

  - **Status 403**

  - Returned if user doesn't have permission to create space backups

  - **Status 409**

  - Returned if backup with the same spaces selected is already in PROGRESS or QUEUED

#### [Find jobs](#api/backup-restore-findJobs)`GET /rest/api/backup-restore/jobs`

Returns jobs based on the filters provided.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `limit` | *int* Default: `25` | amount of jobs that should be returned |
| `fromDate` | *string* | minimum job creation date. Supported date format is \`yyyy-MM-ddTHH:mm:ss.SSSZ\` |
| `toDate` | *string* | maximum job create date. Supported date format is \`yyyy-MM-ddTHH:mm:ss.SSSZ\` |
| `jobScope` | *string* | scope of the job. Acceptable values: "SPACE" and "SITE" |
| `jobOperation` | *string* | job operation. Acceptable values: "BACKUP" and "RESTORE" |
| `jobStates` | *string* | list of job states. Acceptable values: "QUEUED", "PROCESSING", "FINISHED", "CANCELLING", "CANCELLED", "FAILED" |
| `spaceKey` | *string* |  |
| `owner` | *string* | userName of user who had created a job |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the backup/restore jobs 

  - **Status 400**

  - Returned if invalid filter parameters were passed

#### [Get job](#api/backup-restore-getJob)`GET /rest/api/backup-restore/jobs/{jobId}`

Get job by id. The user must be a sysadmin or the owner of the job 

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the backup/restore job 

  - **Status 400**

  - Returned if jobId is null

  - **Status 404**

  - Returned if job not found or user doesn't have permission to see it

#### [Cancel job](#api/backup-restore-cancelJob)`PUT /rest/api/backup-restore/jobs/{jobId}/cancel`

Cancels the job. If the job is already cancelled or failed, the method will do nothing.

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the cancelled job 

  - **Status 404**

  - Returned if job not found or user doesn't have permission to cancel it

#### [Download backup file](#api/backup-restore-downloadBackupFile)`GET /rest/api/backup-restore/jobs/{jobId}/download`

Downloads the backup file for the given job. Requires site admin or space export permissions for all spaces included in the backup job. 

##### Responses

  - **Status 200** - *application/octet-stream*

  - Returns a data stream of the backup file content

  - **Status 404**

  - Returned if job not found or user doesn't have permissions for the job or the file is missing

#### [Cancel all queued jobs](#api/backup-restore-cancelAllQueuedJobs)`PUT /rest/api/backup-restore/jobs/clear-queue`

Cancels all queued jobs. Does not affect jobs that are being processed at the moment.

##### Responses

  - **Status 200** - *application/json*

  - Returns number of cancelled jobs 

  - **Status 403**

  - Returned if user doesn't have permission to cancel jobs

#### [Get files](#api/backup-restore-getFiles)`GET /rest/api/backup-restore/restore/files`

returns list of information on files in conf-home/restore/(jobScope)

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `jobScope` | *string* | name of type of restore job (SITE or SPACE or null), if null, all backup files are listed |

##### Responses

  - *application/json*

#### [Create site restore job](#api/backup-restore-createSiteRestoreJob)`POST /rest/api/backup-restore/restore/site`

Starts the new site restore job.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/site-restore-settings#","title":"Site Restore Settings","type":"object","properties":{"fileName":{"type":"string"},"skipReindex":{"type":"boolean"}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the site restore job 

  - **Status 400**

  - Returned if invalid filename provided

  - **Status 403**

  - Returned if user doesn't have permission to restore site

#### [Create site restore job for uploaded backup file](#api/backup-restore-createSiteRestoreJobForUploadedBackupFile)`POST /rest/api/backup-restore/restore/site/upload`

Starts the new site restore job for the uploaded backup file (Zip). 

This resource expects a multipart post. The media-type multipart/form-data is defined in RFC 1867. Most client libraries have classes that make dealing with multipart posts simple. For instance, in Java the Apache HTTP Components library provides a [MultiPartEntity](http://hc.apache.org/httpcomponents-client-ga/httpmime/apidocs/org/apache/http/entity/mime/MultipartEntity.html) that makes it simple to submit a multipart POST. 

In order to protect against XSRF attacks, because this method accepts multipart/form-data, it has XSRF protection on it. This means you must submit a header of X-Atlassian-Token: nocheck with the request, otherwise it will be blocked. 

The name of the multipart/form-data parameter that contains attachments must be "file" 

An example to attach the file 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.zip" http://myhost/rest/api/backup-restore/restore/site/upload"
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the site restore job 

  - **Status 400**

  - Returned if the uploaded file is not a zip file

  - **Status 403**

  - Returned if user doesn't have permission to restore site

#### [Create space restore job](#api/backup-restore-createSpaceRestoreJob)`POST /rest/api/backup-restore/restore/space`

Creates new space restore job and adds it to the queue.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/space-restore-settings#","title":"Space Restore Settings","type":"object","properties":{"fileName":{"type":"string"},"skipReindex":{"type":"boolean"}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the space restore job 

  - **Status 400**

  - Returned if invalid filename provided

  - **Status 403**

  - Returned if user doesn't have permission to restore spaces

#### [Create space restore job for uploaded backup file](#api/backup-restore-createSpaceRestoreJobForUploadedBackupFile)`POST /rest/api/backup-restore/restore/space/upload`

Starts the new space restore job for the uploaded backup file (Zip). 

This resource expects a multipart post. The media-type multipart/form-data is defined in RFC 1867. Most client libraries have classes that make dealing with multipart posts simple. For instance, in Java the Apache HTTP Components library provides a [MultiPartEntity](http://hc.apache.org/httpcomponents-client-ga/httpmime/apidocs/org/apache/http/entity/mime/MultipartEntity.html) that makes it simple to submit a multipart POST. 

In order to protect against XSRF attacks, because this method accepts multipart/form-data, it has XSRF protection on it. This means you must submit a header of X-Atlassian-Token: nocheck with the request, otherwise it will be blocked. 

The name of the multipart/form-data parameter that contains attachments must be "file" 

An example to attach the file 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.zip" http://myhost/rest/api/backup-restore/restore/space/upload"
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the space restore job 

  - **Status 400**

  - Returned if the uploaded file is not a zip file

  - **Status 403**

  - Returned if user doesn't have permission to restore space

### [api/content](#api/content)

REST wrapper for the ContentService. Provides methods for finding, creating, modifying and deleting Content.Show more

#### [Create content](#api/content-createContent)`POST /rest/api/content`

Creates a new piece of Content or publishes the draft if the content id is present.

For the case publishing draft, a new piece of content will be created and all metadata from the draft will be transferred into the newly created content.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* Default: `current` |  |
| `expand` | *string* Default: `body.storage,history,space,container.history,container.version,version,ancestors` | A comma separated list of properties to expand on the content. Default value: `history,space,version` |

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - *application/json*

#### [Get content](#api/content-getContent)`GET /rest/api/content`

Returns a paginated list of Content.

Example request URI(s):

- `http://example.com/confluence/rest/api/content?spaceKey=TST&title=Cheese&expand=space,body.view,version,container` 

- `http://example.com/confluence/rest/api/content?type=blogpost&spaceKey=TST&title=Bacon&postingDay=2014-02-13&expand=space,body.view,version,container` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `type` | *string* Default: `page` | the content type to return. Default value: `page`. Valid values: `page, blogpost`. |
| `spaceKey` | *string* | the space key to find content under. |
| `title` | *string* | the title of the page to find. Required for `page` type. |
| `status` | *string* | list of statuses the content to be found is in. Defaults to current is not specified. If set to 'any', content in 'current' and 'trashed' status will be fetched. Does not support 'historical' status for now.  |
| `postingDay` | *string* | the posting day of the blog post. Required for `blogpost` type. Format: `yyyy-mm-dd`. Example: `2013-02-13` |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the content. Default value: `history,space,version`. |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `25` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a list of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-content#","title":"Page Response of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if the calling user does not have permission to view the content.

#### [Update](#api/content-update)`PUT /rest/api/content/{contentId}`

Updates a piece of Content, including changes to content status

To update a piece of content you must increment the version.number, supplying the number of the version you are creating. The title property can be updated on all content, body can be updated on all content that has a body (not attachments). For instance to update the content of a blogpost that currently has version 1:

`PUT /rest/api/content/456` 


                                 {

                                 "version":{

                                 "number": 2

                                 },

                                 "title":"My new title",

                                 "type":"page",

                                 "body":{

                                 "storage":{

                                 "value":"<p>New page data.</p>",

                                 "representation":"storage"

                                 }

                                 }

                                 }

                                 
```

To update a page and change its parent page, supply the ancestors property with the request with the parent as the first ancestor i.e. to move a page to be a child of page with ID 789:

`PUT /rest/api/content/456` 


                                 {

                                 "version":{

                                 "number": 2

                                 },

                                 "ancestors": [{"id":789}],

                                 "type":"page",

                                 "body":{

                                 "storage":{

                                 "value":"<p>New page data.</p>",

                                 "representation":"storage"

                                 }

                                 }

                                 }

                                 
```

### Changing status

To restore a piece of content that has the status of trashed the content must have it's version incremented, and status set to current. No other field modifications will be performed when restoring a piece of content from the trash.

Request example to restore from trash: `{"id": "557059","status": "current","version": {"number": 2}}`

If the content you're updating has a draft, specifying status=draft will delete that draft and the body of the content will be replaced with the body specified in the request.

Request example to delete a draft:

`PUT: http://localhost:9096/confluence/rest/api/content/2149384202?status=draft` 


                                 {

                                 "id":"2149384202",

                                 "status":"current",

                                 "version":{

                                 "number":4

                                 },

                                 "space":{

                                 "key":"TST"

                                 },

                                 "type":"page",

                                 "title":"page title",

                                 "body":{

                                 "storage":{

                                 "value":"<p>New page data.</p>",

                                 "representation":"storage"

                                 }

                                 }

                                 }

                                 
```

Updating a draft is not currently supported.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* | the existing status of the content to be updated. |
| `conflictPolicy` | *string* Default: `abort` |  |

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - if no space or no content type, or setup a wrong version type set to content, or status param is not draft and status content is current 

  - **Status 404**

  - if can not find draft with current content

#### [Get content by id](#api/content-getContentById)`GET /rest/api/content/{id}`

Returns a piece of Content.

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234?expand=space,body.view,version,container` 

- `http://example.com/confluence/rest/api/content/1234?status=any` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* | list of Content statuses to filter results on. Default value: `[current]` |
| `version` | *int* |  |
| `expand` | *string* Default: `history,space,version` | A comma separated list of properties to expand on the content. Default value: `history,space,version` We can also specify some extensions such as `extensions.inlineProperties` (for getting inline comment-specific properties) or `extensions.resolution` for the resolution status of each comment in the results  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Delete](#api/content-delete)`DELETE /rest/api/content/{id}`

Trashes or purges a piece of Content, based on its \{@link ContentType\} and \{@link ContentStatus\}. 

There are three cases: 

- If the content is trashable and its status is \{@link ContentStatus\#CURRENT\}, it will be trashed. 

- If the content is trashable, its status is \{@link ContentStatus\#TRASHED\} and the "status" query parameter in the request is "trashed", the content will be purged from the trash and deleted permanently. 

- If the content is not trashable it will be deleted permanently without being trashed. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* | the status of the content to be deleted |

##### Responses

  - **Status 200** - *application/json*

  - Returned if successfully trashed. 

  - **Status 204**

  - Returned if successfully purged.

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to trash or purge the content. 

  - **Status 409**

  - Returned if there is a stale data object conflict when trying to delete a draft

#### [Get history](#api/content-getHistory)`GET /rest/api/content/{id}/history`

Returns the history of a particular piece of content 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/history` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default: `previousVersion,nextVersion,lastUpdated` | the properties on content history to expand |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the content's history 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/history#","title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"$ref":"#/definitions/contributors"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"$ref":"#/definitions/contributors"}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"contributors":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false,"required":["latest"]}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Get macro body by hash](#api/content-getMacroBodyByHash)  deprecated`GET /rest/api/content/{id}/history/{version}/macro/hash/{hash}`

Returns the body of a macro (in storage format) with the given hash. This resource is primarily used by connect applications that require the body of macro to perform their work. 

The hash is generated by connect during render time of the local macro holder and is usually only relevant during the scope of one request. For optimisation purposes, this hash will usually live for multiple requests. 

Collecting a macro by its hash should now be considered deprecated and will be replaced, transparently with macroIds. This resource is currently only called from connect addons which will eventually all use the \{@link \#getContentById(com.atlassian.confluence.api.model.content.id.ContentId, java.util.List, Integer, String)\} resource. 

To make the migration as seamless as possible, this resource will match macros against a generated hash or a stored macroId. This will allow add ons to work during the migration period. 

##### Responses

  - **Status 200** - *application/json*

  - Returns a json representation of a macro. 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/macro-instance#","title":"Macro Instance","type":"object","properties":{"name":{"type":"string"},"body":{"type":"string"},"parameters":{"type":"object","patternProperties":{".+":{"title":"Macro Parameter Instance","type":"object","properties":{"value":{"type":"string"}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content, or there is no macro matching the given hash or id. 

#### [Get macro body by macro id](#api/content-getMacroBodyByMacroId)`GET /rest/api/content/{id}/history/{version}/macro/id/{macroId}`

Returns the body of a macro (in storage format) with the given id. This resource is primarily used by connect applications that require the body of macro to perform their work. 

When content is created, if no macroId is specified, then Confluence will generate a random id. The id is persisted as the content is saved and only modified by Confluence if there are conflicting IDs. 

To preserve backwards compatibility this resource will also match on the hash of the macro body, even if a macroId is found. This check will become redundant as pages get macroId's generated for them and transparently propagate out to all instances. 

##### Responses

  - **Status 200** - *application/json*

  - Returns a json representation of a macro. 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/macro-instance#","title":"Macro Instance","type":"object","properties":{"name":{"type":"string"},"body":{"type":"string"},"parameters":{"type":"object","patternProperties":{".+":{"title":"Macro Parameter Instance","type":"object","properties":{"value":{"type":"string"}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content, or there is no macro matching the given id or hash. 

#### [Scan content](#api/content-scanContent)`GET /rest/api/content/scan`

Returns a paginated list of Content.

Example request URI(s):

- `http://example.com/confluence/rest/api/content/scan?spaceKey=TST&limit=100&expand=space,body.view,version,container` 

- `http://example.com/confluence/rest/api/content/scan?limit=100&expand=space,body.view,version,container` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `spaceKey` | *string* | the space key to find content under. |
| `status` | *string* | list of statuses the content to be found is in. Defaults to current is not specified. If set to 'any', content in 'current', 'draft' and 'trashed' status will be fetched.  |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the content. Default value: `history,space,version`. |
| `limit` | *int* Default: `25` | the limit of the number of items to return, this may be restricted by fixed system limits |
| `cursor` | *string* | the identifier which is used to skip results from a previous query when paginating. Cursor is empty in first request, to move forward or backward use cursor provided in response.  |

##### Responses

  - *application/json*

#### [Search](#api/content-search)`GET /rest/api/content/search`

Fetch a list of content using the Confluence Query Language (CQL). See: [Advanced searching using CQL](https://developer.atlassian.com/display/CONFDEV/Advanced+Searching+using+CQL) 

For example:

Example request URI(s): 

- `http://localhost:8080/confluence/rest/api/content/search?cql=creator=currentUser()&cqlcontext={%22spaceKey%22:%22TST%22, %22contentId%22:%2255%22}` 

- `http://localhost:8080/confluence/rest/api/content/search?cql=space=DEV%20AND%20label=docs&expand=space,metadata.labels&limit=10` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `cql` | *string* | a cql query string to use to locate content |
| `cqlcontext` | *string* | the context to execute a cql search in, this is the json serialized form of SearchContext |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the content. |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `25` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a paginated list of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-content#","title":"Page Response of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 400**

  - Returned if the CQL is invalid or missing

### [api/content/\{id\}/child](#api/content/{id}/child)

REST wrapper for the \{@link ChildContentService\}.

#### [Children](#api/content/{id}/child-children)`GET /rest/api/content/{id}/child`

Returns a map of the direct children of a piece of Content. Content can have multiple types of children - for example a Page can have children that are also Pages, but it can also have Comments and Attachments. 

The \{@link ContentType\}(s) of the children returned is specified by the "expand" query parameter in the request - this parameter can include expands for multiple child types.   
If no types are included in the expand parameter, the map returned will just list the child types that are available to be expanded for the \{@link Content\} referenced by the "id" path parameter. 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child` 

- `http://example.com/confluence/rest/api/content/1234/child?expand=page.body.VIEW` 

- `http://example.com/confluence/rest/api/content/1234/child?expand=page&start=20&limit=10` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the children |
| `parentVersion` | *int* Default: `0` | an int representing the version of the content to retrieve children for |
| `start` | *int* |  |
| `limit` | *int* Default: `25` |  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON map representing multiple ordered collections of content children, keyed by content type 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/map-of-content-type#","title":"Map of Content Type","type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Children of type](#api/content/{id}/child-childrenOfType)`GET /rest/api/content/{id}/child/{type}`

Returns the direct children of a piece of Content, limited to a single child type. 

The \{@link ContentType\}(s) of the children returned is specified by the "type" path parameter in the request. 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/page` 

- `http://example.com/confluence/rest/api/content/1234/child/comment` 

- `http://example.com/confluence/rest/api/content/1234/child/page?expand=body.view` 

- `http://example.com/confluence/rest/api/content/1234/child/comment?start=20&limit=10` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the children |
| `parentVersion` | *int* Default: `0` | an int representing the version of the content to retrieve children for |
| `start` | *int* | (optional, default: 0) the index of the first item within the result set that should be returned |
| `limit` | *int* Default: `25` | (optional, default: site limit) how many items should be returned after the start index |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON map representing multiple ordered collections of content children, keyed by content type 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Comments of content](#api/content/{id}/child-commentsOfContent)`GET /rest/api/content/{id}/child/comment`

Returns the comments of a content 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/comment` 

- `http://example.com/confluence/rest/api/content/1234/child/comment?expand=body.view` 

- `http://example.com/confluence/rest/api/content/1234/child/comment?start=20&limit=10` 

- `http://example.com/confluence/rest/api/content/1234/child/comment?location=footer&location=inline&location=resolved` 

- `http://example.com/confluence/rest/api/content/1234/child/comment?expand=extensions.inlineProperties,extensions.resolution` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the children. We can also specify some extensions such as `extensions.inlineProperties` (for getting inline comment-specific properties) or `extensions.resolution` for the resolution status of each comment in the results  |
| `parentVersion` | *int* Default: `0` | an int representing the version of the content to retrieve children for |
| `start` | *int* | (optional, default: 0) the index of the first item within the result set that should be returned |
| `limit` | *int* Default: `25` | (optional, default: site limit) how many items should be returned after the start index |
| `location` | *string* | (optional, default: "" means all) the location of the comments. Possible values are: "inline", "footer", "resolved". You can define multiple location params. The results will be the comments matched by any location.  |
| `depth` | *string* Default:  | (optional, default: "") the depth of the comments. Possible values are: "" (ROOT only), "all" |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON map representing multiple ordered collections of content children, keyed by content type 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

### [api/content/\{id\}/child/attachment](#api/content/{id}/child/attachment)

CRUD operations for Attachments on Content.

#### [Get attachments](#api/content/{id}/child/attachment-getAttachments)`GET /rest/api/content/{id}/child/attachment`

Returns a paginated list of attachment Content entities within a single container.

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/attachment?start=0&limit=10` 

- `http://example.com/confluence/rest/api/content/1234/child/attachment?filename=myfile.txt&expand=version,container` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the Attachments returned. Optional. |
| `start` | *int* | the index of the first item within the result set that should be returned. Optional. |
| `limit` | *int* Default: `50` | how many items should be returned after the start index. Optional. |
| `filename` | *string* | (optional) filter parameter to return only the Attachment with the matching file name. Optional. |
| `mediaType` | *string* | (optional) filter parameter to return only Attachments with a matching Media-Type. Optional. |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of a list of attachment Content entities 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Create attachments](#api/content/{id}/child/attachment-createAttachments)`POST /rest/api/content/{id}/child/attachment`

Add one or more attachments to a Confluence Content entity, with optional comments.

Comments are optional, but if included there must be as many comments as there are files, and the comments must be in the same order as the files. 

This resource expects a multipart post. The media-type multipart/form-data is defined in RFC 1867. Most client libraries have classes that make dealing with multipart posts simple. For instance, in Java the Apache HTTP Components library provides a [MultiPartEntity](http://hc.apache.org/httpcomponents-client-ga/httpmime/apidocs/org/apache/http/entity/mime/MultipartEntity.html) that makes it simple to submit a multipart POST. 

In order to protect against XSRF attacks, because this method accepts multipart/form-data, it has XSRF protection on it. This means you must submit a header of X-Atlassian-Token: nocheck with the request, otherwise it will be blocked. 

The name of the multipart/form-data parameter that contains attachments must be "file" 

A simple example to attach a file called "myfile.txt" to the container with id "123", with a comment included: 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.txt" -F "comment=This is my File" http://myhost/rest/api/content/123/child/attachment
```

A example to attach a file called "myfile.txt" to the container with id "123", with a comment, and set the minorEdits flag to be true: 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.txt" -F "minorEdit=true" -F "comment=This

                                 is my File" http://myhost/rest/api/content/123/child/attachment
```

An example to attach the same file, with no comment: 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.txt" http://myhost/rest/api/content/123/child/attachment
```

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/attachment` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* Default: `current` | a string containing the status of the attachments content container, supports current or draft, defaults to current |
| `allowDuplicated` | *boolean* Default: `false` | allows to upload an attachment with an existing filename. Optional. |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the Attachments returned. Optional. |

##### Responses

  - **Status 200** - *application/json*

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if attachments is disabled or if you don't have permission to add attachments to this content.

  - **Status 404**

  - Returned if the requested content is not found, the user does not have permission to view it, or if the attachments exceeds the maximum configured attachment size. 

#### [Update](#api/content/{id}/child/attachment-update)`PUT /rest/api/content/{id}/child/attachment/{attachmentId}`

Update the non-binary data of an Attachment. 

This resource can be used to update an attachment's filename, media-type, comment, and parent container. 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/attachment/5678` 

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the attachment id or the attachment version number are invalid.

  - **Status 403**

  - Returned if you are not permitted to update or move the attachment to a different container.

  - **Status 404**

  - Returned if no attachment is found for the attachmentId.

  - **Status 409**

  - Returned if the version of the supplied Attachment does not match the exact version of the Attachment stored in the database.

#### [Update data](#api/content/{id}/child/attachment-updateData)`POST /rest/api/content/{id}/child/attachment/{attachmentId}/data`

Update the binary data of an Attachment, and optionally the comment and the minor edit field.

This adds a new version of the attachment, containing the new binary data, filename, and content-type.

When updating the binary data of an attachment, the comment related to it together with the field that specifies if it's a minor edit can be updated as well, but are not required. If an update is considered to be a minor edit, notifications will not be sent to the watchers of that content. 

This resource expects a multipart post. The media-type multipart/form-data is defined in RFC 1867. Most client libraries have classes that make dealing with multipart posts simple. For instance, in Java the Apache HTTP Components library provides a [MultiPartEntity](http://hc.apache.org/httpcomponents-client-ga/httpmime/apidocs/org/apache/http/entity/mime/MultipartEntity.html) that makes it simple to submit a multipart POST. 

In order to protect against XSRF attacks, because this method accepts multipart/form-data, it has XSRF protection on it. This means you must submit a header of X-Atlassian-Token: nocheck with the request, otherwise it will be blocked. 

The name of the multipart/form-data parameter that contains attachments must be "file" 

A simple example to upload a file called "myfile.txt" to the Attachment with id "456" in a container with id "123", with the comment updated, and minorEdit set to true: 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.txt" -F "minorEdit=true" -F "comment=This

                                 is my updated File" http://myhost/rest/api/content/123/child/attachment/456/data
```

An example to upload the same file, with no comment: 

```
curl -D- -u admin:admin -X POST -H "X-Atlassian-Token: nocheck" -F "file=@myfile.txt" http://myhost/rest/api/content/123/child/attachment/456/data
```

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/child/attachment/5678/data` 

##### Responses

  - **Status 200** - *application/json*

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the attachment id is invalid.

  - **Status 404**

  - Returned if no attachment is found for the attachmentId.

### [api/content/\{id\}/descendant](#api/content/{id}/descendant)

REST wrapper for the \{@link ChildContentService\}, when \{@link Depth\} is ALL.

#### [Descendants](#api/content/{id}/descendant-descendants)`GET /rest/api/content/{id}/descendant`

Returns a map of the descendants of a piece of Content. Content can have multiple types of descendants - for example a Page can have descendants that are also Pages, but it can also have Comments and Attachments. 

The \{@link ContentType\}(s) of the descendants returned is specified by the "expand" query parameter in the request - this parameter can include expands for multiple descendant types.   
If no types are included in the expand parameter, the map returned will just list the descendant types that are available to be expanded for the \{@link Content\} referenced by the "id" path parameter. 

Currently the only supported descendants are comment descendants of non-comment Content. 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/descendant` 

- `http://example.com/confluence/rest/api/content/1234/descendant?expand=comment.body.VIEW` 

- `http://example.com/confluence/rest/api/content/1234/descendant?expand=comment` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the descendants |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON map representing multiple ordered collections of content descendants, keyed by content type 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/map-of-content-type#","title":"Map of Content Type","type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Descendants of type](#api/content/{id}/descendant-descendantsOfType)`GET /rest/api/content/{id}/descendant/{type}`

Returns the direct descendants of a piece of Content, limited to a single descendant type. 

The \{@link ContentType\}(s) of the descendants returned is specified by the "type" path parameter in the request. 

Currently the only supported descendants are comment descendants of non-comment Content. 

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/descendant/comment` 

- `http://example.com/confluence/rest/api/content/1234/descendant/comment?expand=body.VIEW` 

- `http://example.com/confluence/rest/api/content/1234/descendant/comment?start=20&limit=10` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the descendants |
| `start` | *int* | (optional, default: 0) the index of the first item within the result set that should be returned |
| `limit` | *int* Default: `25` | (optional, default: site limit) how many items should be returned after the start index |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON map representing multiple ordered collections of content descendants, keyed by content type 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

### [api/content/\{id\}/label](#api/content/{id}/label)

REST wrapper for the \{@link ContentLabelService\}.

#### [Add labels](#api/content/{id}/label-addLabels)`POST /rest/api/content/{id}/label`

Adds a list of labels to the specified content.

The body is the json representation of the list.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/string#","title":"String","type":"string"}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the labels on a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-label#","title":"Page Response of Label","type":"array","items":{"title":"Label","type":"object","properties":{"prefix":{"type":"string"},"name":{"type":"string"},"id":{"type":"string"}},"additionalProperties":false}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Delete label with query param](#api/content/{id}/label-deleteLabelWithQueryParam)`DELETE /rest/api/content/{id}/label`

Deletes a labels to the specified content.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `name` | *string* | the name of the label to be removed from the content |

##### Responses

  - **Status 204**

  - Empty response on successful delete

  - **Status 403**

  - Returned if user has view permission, but no edit permission to the content

  - **Status 404**

  - Returned if content or label doesn't exist, or calling user doesn't have view permission to the content

#### [Labels](#api/content/{id}/label-labels)`GET /rest/api/content/{id}/label`

Returns the list of labels on a piece of Content.

Example request URI(s):

- `http://example.com/confluence/rest/api/content/1234/label` 

- `http://example.com/confluence/rest/api/content/1234/label?prefix=global&start=0&limit=200` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `prefix` | *string* | the prefixes to filter the labels with \{@see Label.Prefix\} |
| `start` | *int* |  |
| `limit` | *int* Default: `200` | the limit of the number of labels to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the labels on a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-label#","title":"Page Response of Label","type":"array","items":{"title":"Label","type":"object","properties":{"prefix":{"type":"string"},"name":{"type":"string"},"id":{"type":"string"}},"additionalProperties":false}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Delete label](#api/content/{id}/label-deleteLabel)`DELETE /rest/api/content/{id}/label/{label}`

Deletes a labels to the specified content. When calling this method through REST the label parameter doesn't accept "/" characters in label names, because of security constraints. For this case please use the query parameter version of this method (/content/\{id\}/label?name=\{label\}

##### Responses

  - **Status 204**

  - Empty response on successful delete

  - **Status 400**

  - Returned if trying to delete a label with "/" in the name

  - **Status 403**

  - Returned if user has view permission, but no edit permission to the content

  - **Status 404**

  - Returned if content or label doesn't exist, or calling user doesn't have view permission to the content

### [api/content/\{id\}/property](#api/content/{id}/property)

A REST resource for manipulating content properties.

Content properties are a key / value store of properties attached to a piece of Content. The key is a string, and the value is a JSON object.

Show more

#### [Find all](#api/content/{id}/property-findAll)`GET /rest/api/content/{id}/property`

Returns a paginated list of content properties.

Example request URI:

- `http://example.com/confluence/rest/api/content/1234/property?expand=content,version` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* | a comma separated list of properties to expand on the content properties. Default value: `version`. |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `10` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the content property list 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-json-content-property#","title":"Rest List of Json Content Property","type":"array","items":{"title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Create](#api/content/{id}/property-create)`POST /rest/api/content/{id}/property`

Creates a new content property.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the content property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the given property has a different ContentId to the one in the path, or if the content already has a value with the given key, or the value is missing, or the value is too long. 

  - **Status 403**

  - Returned if the user does not have permission to edit the content with the given ContentId

  - **Status 413**

  - Returned if the value is too long.

#### [Find by key](#api/content/{id}/property-findByKey)`GET /rest/api/content/{id}/property/{key}`

Returns a content property.

Example request URI:

- `http://example.com/confluence/rest/api/content/1234/property/example-property-key?expand=content,version` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* | a comma separated list of properties to expand on the content properties. Default value: `version`. |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the content property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or no property with the given key, or if the calling user does not have permission to view the content. 

#### [Update](#api/content/{id}/property-update)`PUT /rest/api/content/{id}/property/{key}`

Updates a content property.

The body contains the representation of the content property. Must include the property id, and the new version number. Attempts to create a new content property if the given version number is 1, just like \{@link \#create(com.atlassian.confluence.api.model.content.id.ContentId, String, com.atlassian.confluence.api.model.content.JsonContentProperty)\}.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the given property has a different ContentId to the one in the path, or if the property has a different key to the one in the path, or the value is missing, or the value is too long. 

  - **Status 403**

  - Returned if the user does not have permission to edit the content with the given ContentId

  - **Status 404**

  - Returned if there is no content with the given id, or no property with the given key, or if the calling user does not have permission to view the content. 

  - **Status 409**

  - Returned if the given version is does not match the expected target version of the updated property

  - **Status 413**

  - Returned if the value is too long.

#### [Delete](#api/content/{id}/property-delete)`DELETE /rest/api/content/{id}/property/{key}`

Deletes a content property.

##### Responses

  - **Status 204**

  - Returned if successfully deleted.

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Create](#api/content/{id}/property-create)`POST /rest/api/content/{id}/property/{key}`

Creates a new content property.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the content property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-content-property#","title":"Json Content Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"id":{"title":"Json Content Property Id","type":"object"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the given property has a different ContentId to the one in the path, or if the content already has a value with the given key, or the value is missing, or the value is too long. 

  - **Status 403**

  - Returned if the user does not have permission to edit the content with the given ContentId

  - **Status 413**

  - Returned if the value is too long.

### [api/content/\{id\}/restriction](#api/content/{id}/restriction)

#### [By operation](#api/content/{id}/restriction-byOperation)`GET /rest/api/content/{id}/restriction/byOperation`

Returns info about all restrictions by operation

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default: `update.restrictions.user,read.restrictions.group,read.restrictions.user,update.restrictions.group` | a comma separated list of properties to expand on the content properties. Default value: group. |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the restrictions group by operations 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/map-of-operation-key-and-content-restriction#","title":"Map of Operation Key-and-Content Restriction","type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

#### [For operation](#api/content/{id}/restriction-forOperation)`GET /rest/api/content/{id}/restriction/byOperation/{operationKey}`

Returns info about all restrictions of given operation

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default: `restrictions.user,restrictions.group` | a comma separated list of properties to expand on the content properties. Default value: group. |
| `start` | *int* | pagination start |
| `limit` | *int* Default: `100` | pagination limit |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation of the restrictions of given operation 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content-restriction#","title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"$ref":"#/definitions/subject-restrictions"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"$ref":"#/definitions/subject-restrictions"}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"subject-restrictions":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

### [api/content/blueprint](#api/content/blueprint)

#### [Publish legacy draft](#api/content/blueprint-publishLegacyDraft)`POST /rest/api/content/blueprint/instance/{draftId}`

Publishes a legacy draft of a Content created from a ContentBlueprint

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* Default: `draft` | only support 'draft' status for now. |
| `expand` | *string* Default: `body.storage,history,space,version,ancestors` | A comma separated list of properties to expand on the content. Default value: `body.storage,history,space,version,ancestors` |

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - *application/json*

#### [Publish shared draft](#api/content/blueprint-publishSharedDraft)`PUT /rest/api/content/blueprint/instance/{draftId}`

Publishes a shared draft of a Content created from a ContentBlueprint

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `status` | *string* Default: `draft` | only support 'draft' status for now. |
| `expand` | *string* Default: `body.storage,history,space,version,ancestors` | A comma separated list of properties to expand on the content. Default value: `body.storage,history,space,version,ancestors` |

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content#","title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"container":{"title":"Container","type":"object"},"content":{"title":"Content","type":"object","properties":{"id":{"$ref":"#/definitions/content-id"},"type":{"$ref":"#/definitions/content-type"},"status":{"$ref":"#/definitions/content-status"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"$ref":"#/definitions/history"}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"$ref":"#/definitions/operation-check-result"}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"$ref":"#/definitions/container"}},"body":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-body"}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/content-restriction"}},"additionalProperties":false}},"additionalProperties":false},"content-body":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"content-id":{"title":"Content Id","type":"object"},"content-representation":{"title":"Content Representation","type":"object"},"content-restriction":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false},"content-status":{"title":"Content Status","type":"object"},"content-type":{"title":"Content Type","type":"object"},"history":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-check-result":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - *application/json*

### [api/contentbody/convert/\{to\}](#api/contentbody/convert/{to})

Used for converting ContentBody objects from one format to another.

#### [Convert](#api/contentbody/convert/{to}-convert)`POST /rest/api/contentbody/convert/{to}`

Converts between content body representations.

Not all representations can be converted to/from other formats. Supported conversions:

| **Source Representation** | **Destination Representation Supported** |
| :-: | :-: |
| storage | view,export\_view,styled\_view,editor |
| editor | storage |
| view | *None* |
| export\_view | *None* |
| styled\_view | *None* |

Example request URI(s):

- `http://example.com/confluence/rest/api/contentbody/convert/view` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  |  |

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/content-body#","title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content-body#","title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

### [api/group](#api/group)

Non-admin operations for user groups

#### [Get groups](#api/group-getGroups)`GET /rest/api/group`

Get a paginated collection of user groups

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  |  |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `200` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a list of groups 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-group#","title":"Page Response of Group","type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view groups.

#### [Get group](#api/group-getGroup)`GET /rest/api/group/{groupName}`

Get the user group with the group name

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* |  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a single groups 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/group#","title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view groups.

#### [Get members](#api/group-getMembers)`GET /rest/api/group/{groupName}/member`

**You need admin permissions to use this endpoint.** Get a paginated collection of users in the given group 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* |  |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `200` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a single groups 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-person#","title":"Page Response of Person","type":"array","items":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view users.

### [api/longtask](#api/longtask)

REST wrapper for the LongTaskService.

#### [Get tasks](#api/longtask-getTasks)`GET /rest/api/longtask`

Returns information about all tracked long-running tasks.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the tasks |
| `start` | *int* |  |
| `limit` | *int* Default: `100` |  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a list of long tasks 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-long-task-status#","title":"Rest List of Long Task Status","type":"array","items":{"title":"Long Task Status","type":"object","properties":{"id":{"title":"Long Task Id","type":"object"},"name":{"$ref":"#/definitions/message"},"elapsedTime":{"type":"integer"},"percentageComplete":{"type":"integer"},"successful":{"type":"boolean"},"messages":{"type":"array","items":{"$ref":"#/definitions/message"}}},"additionalProperties":false,"required":["elapsedTime","percentageComplete","successful"]},"definitions":{"message":{"title":"Message","type":"object","properties":{"args":{"type":"array","items":{}},"key":{"type":"string"}},"additionalProperties":false}}}
```

#### [Get task](#api/longtask-getTask)`GET /rest/api/longtask/{id}`

Returns information about a long-running task.

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the task |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a long task 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/long-task-status#","title":"Long Task Status","type":"object","properties":{"id":{"title":"Long Task Id","type":"object"},"name":{"$ref":"#/definitions/message"},"elapsedTime":{"type":"integer"},"percentageComplete":{"type":"integer"},"successful":{"type":"boolean"},"messages":{"type":"array","items":{"$ref":"#/definitions/message"}}},"definitions":{"message":{"title":"Message","type":"object","properties":{"args":{"type":"array","items":{}},"key":{"type":"string"}},"additionalProperties":false}},"additionalProperties":false,"required":["elapsedTime","percentageComplete","successful"]}
```

  - **Status 404**

  - Returned if there is no task with the given key, or if the calling user does not have permission to view it

### [api/search](#api/search)

API REST Resource for searching for entities in Confluence. Generally delegates to the CQLSearchService.Show more

#### [Search](#api/search-search)`GET /rest/api/search`

Search for entities in Confluence using the [Confluence Query Language (CQL)](https://developer.atlassian.com/confdev/confluence-rest-api/advanced-searching-using-cql) 

For example:

Example request URI(s): 

- `http://localhost:8080/confluence/rest/api/search?cql=creator=currentUser()&type%20in%20(space,page,user)&cqlcontext={%22spaceKey%22:%22TST%22, %22contentId%22:%2255%22}` 

- `http://localhost:8080/confluence/rest/api/search?cql=siteSearch~'example'%20AND%20label=docs&expand=content.space,space.homepage&limit=10` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `cql` | *string* | the CQL query see [advanced searching in confluence using CQL](https://developer.atlassian.com/confdev/confluence-rest-api/advanced-searching-using-cql)  |
| `cqlcontext` | *string* | the execution context for CQL functions, provides current space key and content id. If this is not provided some CQL functions will not be available.  |
| `excerpt` | *string* Default: `highlight` | the excerpt strategy to apply to the result, one of : indexed, highlight, none. This defaults to highlight. |
| `expand` | *string* Default:  | the properties to expand on the search result, this may cause database requests for some properties |
| `start` | *int* |  |
| `limit` | *int* Default: `25` | the limit of the number of items to return, this may be restricted by fixed system limits |
| `includeArchivedSpaces` | *boolean* Default: `false` | whether to include content in archived spaces in the result, this defaults to false |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a list of search results 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/search-page-response-of-search-result#","title":"Search Page Response of Search Result","type":"array","items":{"title":"Search Result","type":"object","properties":{"title":{"type":"string"},"excerpt":{"type":"string"},"url":{"type":"string"},"resultParentContainer":{"type":"array","items":{"$ref":"#/definitions/container-summary"}},"resultGlobalContainer":{"type":"array","items":{"$ref":"#/definitions/container-summary"}},"iconCssClass":{"type":"string"},"lastModified":{"type":"string"},"friendlyLastModified":{"type":"string"}},"additionalProperties":false},"definitions":{"container-summary":{"title":"Container Summary","type":"object","properties":{"title":{"type":"string"},"displayUrl":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 400**

  - Returned if the query cannot be parsed

### [api/space](#api/space)

REST wrapper for the SpaceService.

#### [Create space](#api/space-createSpace)`POST /rest/api/space`

Creates a new Space. 

The incoming Space does not include an id, but must include a Key and Name, and should include a Description. 

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a space 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

#### [Spaces](#api/space-spaces)`GET /rest/api/space`

Returns information about a number of spaces.

Example request URI(s):

- `http://example.com/confluence/rest/api/space?spaceKey=TST&spaceKey=ds` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `spaceKey` | *string* | a list of space keys |
| `type` | *string* | filter the list of spaces returned by type (global, personal) |
| `status` | *string* | filter the list of spaces returned by status (current, archived) |
| `label` | *string* | filter the list of spaces returned by label |
| `favourite` | *boolean* | filter the list of spaces returned by favourites |
| `hasRetentionPolicy` | *boolean* | filter the list of spaces returned by retention policy |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the spaces |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `25` | the limit of the number of spaces to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns an array of full JSON representations of found spaces 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-space#","title":"Page Response of Space","type":"array","items":{"$ref":"#/definitions/space"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

#### [Create private space](#api/space-createPrivateSpace)`POST /rest/api/space/_private`

Creates a new private Space, viewable only by its creator. 

The incoming Space does not include an id, but must include a Key and Name, and should include a Description. 

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a space 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

#### [Update](#api/space-update)`PUT /rest/api/space/{spaceKey}`

Updates a Space.

Currently only the Space name, description and homepage can be updated.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a space 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no space with the given key, or if the calling user does not have permission to update it. 

#### [Delete](#api/space-delete)`DELETE /rest/api/space/{spaceKey}`

Deletes a Space.

The space is deleted in a long running task, so the space cannot be considered deleted when this resource returns. Clients can follow the status link in the response and poll it until the task completes. 

##### Responses

  - **Status 202** - *application/json*

  - Returns a pointer to the status of the space-deletion task 

  - **Status 404**

  - Returned if there is no space with the given key, or if the calling user does not have permission to delete it. 

#### [Space](#api/space-space)`GET /rest/api/space/{spaceKey}`

Returns information about a space.

Example request URI(s):

- `http://example.com/confluence/rest/api/space/TST?expand=description` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the space |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a space 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/space#","title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/formatted-body"}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"$ref":"#/definitions/space-type"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"$ref":"#/definitions/space-retention-policy"}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"formatted-body":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space-retention-policy":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]},"space-type":{"title":"Space Type","type":"object"},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no space with the given key, or if the calling user does not have permission to view the space. 

#### [Contents](#api/space-contents)`GET /rest/api/space/{spaceKey}/content`

Returns the content in this given space

Example request URI(s):

- `http://example.com/confluence/rest/api/space/TEST/content?expand=history` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `depth` | *string* Default: `all` | a string indicating if all content, or just the root content of the space is returned. Default value: `all`. Valid values: `all, root`.  |
| `expand` | *string* Default:  | a comma separated list of properties to expand on each piece of content retrieved |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `25` | the limit of the number of labels to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/map-of-content-type#","title":"Map of Content Type","type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

#### [Contents with type](#api/space-contentsWithType)`GET /rest/api/space/{spaceKey}/content/{type}`

Returns the content in this given space with the given type

Example request URI(s):

- `http://example.com/confluence/rest/api/space/TEST/content/page?expand=history` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `depth` | *string* Default: `all` | a string indicating if all content, or just the root content of the space is returned. Default value: `all`. Valid values: `all, root`.  |
| `expand` | *string* Default:  | a comma separated list of properties to expand on each piece of content retrieved |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `25` | the limit of the number of labels to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a piece of content 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-list-of-content#","title":"Rest List of Content","type":"array","items":{"$ref":"#/definitions/content"},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no content with the given id, or if the calling user does not have permission to view the content. 

### [api/space/\{spaceKey\}/property](#api/space/{spaceKey}/property)

A REST resource for manipulating space properties.

Space properties are a key / value store of properties attached to a Space. The key is a string, and the value is a JSON object.

Show more

#### [Get](#api/space/{spaceKey}/property-get)`GET /rest/api/space/{spaceKey}/property`

Returns a paginated list of space properties.

Example request URI:

- `http://example.com/confluence/rest/api/space/TST/property?expand=space,version` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default: `version` | a comma separated list of properties to expand on the space properties. Default value: `version`. |
| `start` | *int* | the start point of the collection to return |
| `limit` | *int* Default: `10` | the limit of the number of items to return, this may be restricted by fixed system limits |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the space property list 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-json-space-property#","title":"Page Response of Json Space Property","type":"array","items":{"title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"additionalProperties":false},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}}}
```

  - **Status 404**

  - Returned if there is no space with the given key, or if the calling user does not have permission to view the space. 

#### [Create](#api/space/{spaceKey}/property-create)`POST /rest/api/space/{spaceKey}/property`

Creates a new space property.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the space property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the space already has a value with the given key, or no property value was provided, or the value is too long.

  - **Status 403**

  - Returned if the user does not have permission to edit the space with the given key

  - **Status 413**

  - Returned if the value is too long.

#### [Get](#api/space/{spaceKey}/property-get)`GET /rest/api/space/{spaceKey}/property/{key}`

Returns a paginated list of space properties.

Example request URI:

- `http://example.com/confluence/rest/api/space/TST/property?expand=space,version` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default: `version` | a comma separated list of properties to expand on the space properties. Default value: `version`. |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the space property list 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if there is no space with the given key, or if the calling user does not have permission to view the space. 

#### [Update](#api/space/{spaceKey}/property-update)`PUT /rest/api/space/{spaceKey}/property/{key}`

Updates a space property.

The body contains the representation of the space property. Must include new version number. If the given version number is 1, attempts to create a new space property, just like \{@link \#create(String, com.atlassian.confluence.api.model.content.JsonSpaceProperty)\}.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the given property has a different spaceKey to the one in the path, or if the property has a different key to the one in the path, or no property value was provided, or the value is too long. 

  - **Status 403**

  - Returned if the user does not have permission to edit the space with the given spaceKey

  - **Status 404**

  - Returned if there is no space with the given spaceKey, or no property with the given key, or if the calling user does not have permission to view the space. 

  - **Status 409**

  - Returned if the given version is does not match the expected target version of the updated property

  - **Status 413**

  - Returned if the value is too long.

#### [Delete](#api/space/{spaceKey}/property-delete)`DELETE /rest/api/space/{spaceKey}/property/{key}`

Deletes a space property.

##### Responses

  - **Status 204**

  - Returned if successfully deleted.

  - **Status 404**

  - Returned if there is no space with the give key, property with the given property key, or if the calling user does not have permission to view the space. 

#### [Create](#api/space/{spaceKey}/property-create)`POST /rest/api/space/{spaceKey}/property/{key}`

Creates a new space property.

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of the space property 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/json-space-property#","title":"Json Space Property","type":"object","properties":{"key":{"type":"string"},"value":{"title":"Json String","type":"object"},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"space":{"type":"array","items":{"$ref":"#/definitions/space"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 400**

  - Returned if the space already has a value with the given key, or no property value was provided, or the value is too long.

  - **Status 403**

  - Returned if the user does not have permission to edit the space with the given key

  - **Status 413**

  - Returned if the value is too long.

### [api/user](#api/user)

Non-admin user operations

#### [Get user](#api/user-getUser)`GET /rest/api/user`

Get information about a user identified by either user key or username.

Example request URI(s):

- `http://example.com/confluence/rest/api/user?username=jblogs` 

- `http://example.com/confluence/rest/api/user?key=402880824ff933a4014ff9345d7c0002` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to request from this resource |
| `username` | *string* | username of the user to request from this resource |
| `expand` | *string* | properties to expand on the user |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a user 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/person#","title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}],"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view users.

  - **Status 404**

  - Returned if a user with the given username or userkey does not exist

#### [Get anonymous](#api/user-getAnonymous)`GET /rest/api/user/anonymous`

Get information about the how anonymous is represented in confluence 

Example request URI(s):

- `http://example.com/confluence/rest/api/user/anonymous` 

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a user 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/person#","title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}],"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to use confluence.

#### [Get current](#api/user-getCurrent)`GET /rest/api/user/current`

Get information about the current logged in user. 

Example request URI(s):

- `http://example.com/confluence/rest/api/user/current` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* |  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a user 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/person#","title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}],"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to use confluence.

#### [Change password](#api/user-changePassword)`POST /rest/api/user/current/password`

Change the password for the current user 

Validation Rules : New password supplied cannot be null or blank 

##### Request

###### Example

```
http://yourhost:8080/confluence/rest/api/user/current/password
```

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/password-change-details#","title":"Password Change Details","type":"object","properties":{"oldPassword":{"type":"string"},"newPassword":{"type":"string"}},"additionalProperties":false}
```

##### Responses

  - **Status 204**

  - returned if the password changed successfully

  - **Status 400**

  - returned if any error occurs while changing user password. Refer the validation rules above

  - **Status 401**

  - returned if the user is not authenticated

  - **Status 403**

  - returned if current password is wrong or if the user has exceeded number of allowed failed login attempts

#### [Get users](#api/user-getUsers)`GET /rest/api/user/list`

Gets a paginated collection of all registered users, including but not limited to: 

- Disabled users 

- Enabled users 

- Enabled users which count towards the license count on the site 

- Enabled users which do not count towards the license count on the site 

- Enabled users which have "can use" global permissions 

- Enabled users which do not have "can use" global permissions 

Example request URI(s): 

- `http://example.com/confluence/rest/api/user/list` 

- `http://example.com/confluence/rest/api/user/list?start=0` 

- `http://example.com/confluence/rest/api/user/list?start=0&limit=100` 

- `http://example.com/confluence/rest/api/user/list?start=0&limit=100&expand=status` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `expand` | *string* Default:  | a comma separated list of properties to expand on the users. Default value is empty string. |
| `start` | *int* Default: `0` | the start point of the collection to return. This must be non-negative. Default value is 0. |
| `limit` | *int* Default: `100` | the limit of the number of users to return, this may be restricted by fixed system limits. This must be non-negative. Default value is 200.  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a collection of users. 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-person#","title":"Page Response of Person","type":"array","items":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false}}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view users. This is possible for anonymous or un-licensed users.

#### [Get groups](#api/user-getGroups)`GET /rest/api/user/memberof`

Get a paginated collection of groups that the given user is a member of 

Example request URI(s):

- `http://example.com/confluence/rest/api/user/memberof?username=jblogs` 

- `http://example.com/confluence/rest/api/user/memberof?key=402880824ff933a4014ff9345d7c0002` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* |  |
| `username` | *string* |  |
| `expand` | *string* |  |
| `start` | *int* |  |
| `limit` | *int* Default: `200` |  |

##### Responses

  - **Status 200** - *application/json*

  - Returns a full JSON representation of a collection of groups 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/page-response-of-group#","title":"Page Response of Group","type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}
```

  - **Status 403**

  - Returned if the calling user does not have permission to use confluence.

### [api/user/\{username\}/group/\{groupName\}](#api/user/{username}/group/{groupName})

REST API for managing \{@link User\} to \{@link Group\} membership.

#### [Update](#api/user/{username}/group/{groupName}-update)`PUT /rest/api/user/{username}/group/{groupName}`

Add the given \{@link User\} identified by username to the given \{@link Group\} identified by groupName. This method is idempotent i.e. if the membership already exists then no action will be taken. 

##### Responses

  - *application/json*

#### [Delete](#api/user/{username}/group/{groupName}-delete)`DELETE /rest/api/user/{username}/group/{groupName}`

Remove the given \{@link User\} identified by username from the given \{@link Group\} identified by groupName. This method is idempotent i.e. if the membership is not present then no action will be taken. 

##### Responses

  - *application/json*

### [api/user/watch](#api/user/watch)

A REST resource for manipulating watchers for \{@link com.atlassian.confluence.api.model.content.Space Spaces\}, pieces of \{@link com.atlassian.confluence.api.model.content.Content Content\} and \{@link com.atlassian.confluence.api.model.content.Label Labels\}.Show more

#### [Add content watcher](#api/user/watch-addContentWatcher)`POST /rest/api/user/watch/content/{contentId}`

Create a new watcher for the given user and content id. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `POST http://example.com/confluence/rest/api/user/watch/content/131213` 

- `POST http://example.com/confluence/rest/api/user/watch/content/131213?username=jblogs` 

- `POST http://example.com/confluence/rest/api/user/watch/content/131213?key=ff8080815a58e24c015a58e263710000` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to create the new watcher for |
| `username` | *string* | username of the user to create the new watcher for |

##### Responses

  - **Status 204** - *application/json*

  - Returned if the watcher was successfully created 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/content-watch#","title":"Content Watch","type":"object","properties":{"watcher":{"$ref":"#/definitions/user"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"title":"Content Type","type":"object"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if no content exists for the specified content id or the calling user does not have permission to perform the operation

#### [Remove content watcher](#api/user/watch-removeContentWatcher)`DELETE /rest/api/user/watch/content/{contentId}`

Delete an existing watcher for the given user and content id. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `DELETE http://example.com/confluence/rest/api/user/watch/content/131213` 

- `DELETE http://example.com/confluence/rest/api/user/watch/content/131213?username=jblogs` 

- `DELETE http://example.com/confluence/rest/api/user/watch/content/131213?key=ff8080815a58e24c015a58e263710000` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to delete the watcher for |
| `username` | *string* | username of the user to delete the watcher for |

##### Responses

  - **Status 204** - *application/json*

  - Returned if the watcher was successfully deleted 

  - **Status 404**

  - Returned if no content exists for the specified content id or the calling user does not have permission to perform the operation

#### [Is watching content](#api/user/watch-isWatchingContent)`GET /rest/api/user/watch/content/{contentId}`

Get information about whether a user is watching a specified content. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `http://example.com/confluence/rest/api/user/watch/content/131213` 

- `http://example.com/confluence/rest/api/user/watch/content/131213?username=jblogs` 

- `http://example.com/confluence/rest/api/user/watch/content/131213?key=ff8080815a58e24c015a58e263710000` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to check for watching state |
| `username` | *string* | username of the user to check for watching state |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation containing the watching state 

  - **Status 404**

  - Returned if no content exists for the specified content id or calling user does not have permission to perform the operation

#### [Add space watch](#api/user/watch-addSpaceWatch)`POST /rest/api/user/watch/space/{spaceKey}`

Create a new watcher for the given user and space key. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `POST http://example.com/confluence/rest/api/user/watch/space/SPACEKEY` 

- `POST http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?username=jblogs` 

- `POST http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?key=ff8080815a58e24c015a58e263710000` 

- `POST http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?contentType=blogpost` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userKey of the user to create the new watcher for |
| `username` | *string* | userName of the user to create the new watcher for |
| `contentType` | *string* | the optional content type to delete the watcher for |

##### Responses

  - **Status 200** - *application/json*

  - Returned if the watcher was successfully created 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/space-watch#","title":"Space Watch","type":"object","properties":{"watcher":{"$ref":"#/definitions/user"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"contentTypes":{"type":"array","items":{"$ref":"#/definitions/content-type"}}},"definitions":{"anonymous":{"title":"Anonymous","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"type":{"type":"string"}},"additionalProperties":false},"content":{"title":"Content","type":"object","properties":{"id":{"title":"Content Id","type":"object"},"type":{"$ref":"#/definitions/content-type"},"status":{"title":"Content Status","type":"object"},"title":{"type":"string"},"space":{"type":"array","items":{"$ref":"#/definitions/space"}},"history":{"type":"array","items":{"title":"History","type":"object","properties":{"previousVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"nextVersion":{"type":"array","items":{"$ref":"#/definitions/version"}},"lastUpdated":{"type":"array","items":{"$ref":"#/definitions/version"}},"latest":{"type":"boolean"},"createdBy":{"$ref":"#/definitions/person"},"createdDate":{"type":"string"},"contributors":{"type":"array","items":{"title":"Contributors","type":"object","properties":{"publishers":{"type":"array","items":{"title":"Contributor Users","type":"object","properties":{"users":{"type":"array","items":{"$ref":"#/definitions/person"}},"userKeys":{"type":"array","items":{"type":"string"}}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false,"required":["latest"]}},"version":{"type":"array","items":{"$ref":"#/definitions/version"}},"ancestors":{"type":"array","items":{"$ref":"#/definitions/content"}},"operations":{"type":"array","items":{"title":"Operation Check Result","type":"object","properties":{"operation":{"$ref":"#/definitions/operation-key"}},"additionalProperties":false}},"children":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"descendants":{"type":"object","patternProperties":{".+":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false},"container":{"type":"array","items":{"title":"Container","type":"object"}},"body":{"type":"object","patternProperties":{".+":{"title":"Content Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false}},"additionalProperties":false},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"extensions":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"restrictions":{"type":"object","patternProperties":{".+":{"title":"Content Restriction","type":"object","properties":{"content":{"type":"array","items":{"$ref":"#/definitions/content"}},"operation":{"$ref":"#/definitions/operation-key"},"restrictions":{"type":"array","items":{"title":"Subject Restrictions","type":"object","properties":{"user":{"type":"array","items":{"$ref":"#/definitions/user"}},"group":{"type":"array","items":{"title":"Group","type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}},"additionalProperties":false},"content-representation":{"title":"Content Representation","type":"object"},"content-type":{"title":"Content Type","type":"object"},"html-string":{"title":"Html String","type":"object"},"icon":{"title":"Icon","type":"object","properties":{"path":{"type":"string"},"width":{"type":"integer"},"height":{"type":"integer"},"isDefault":{"type":"boolean"}},"additionalProperties":false,"required":["width","height","isDefault"]},"known-user":{"title":"Known User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"},"status":{"type":"array","items":{"title":"User Status","type":"object"}}},"additionalProperties":false},"operation-key":{"title":"Operation Key","type":"object"},"person":{"title":"Person","type":"object","anyOf":[{"$ref":"#/definitions/anonymous"},{"$ref":"#/definitions/known-user"},{"$ref":"#/definitions/unknown-user"},{"$ref":"#/definitions/user"}]},"retention-rule":{"title":"Retention Rule","type":"object","properties":{"maxAge":{"type":"integer"},"ageUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"maxNumberOfVersions":{"type":"integer"},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"space":{"title":"Space","type":"object","properties":{"id":{"type":"integer"},"key":{"type":"string"},"name":{"type":"string"},"icon":{"type":"array","items":{"$ref":"#/definitions/icon"}},"description":{"type":"object","patternProperties":{".+":{"title":"Formatted Body","type":"object","properties":{"value":{"type":"string"},"webresource":{"type":"array","items":{"$ref":"#/definitions/web-resource-dependencies"}},"representation":{"$ref":"#/definitions/content-representation"}},"additionalProperties":false}},"additionalProperties":false},"homepage":{"type":"array","items":{"$ref":"#/definitions/content"}},"type":{"title":"Space Type","type":"object"},"metadata":{"type":"object","patternProperties":{".+":{}},"additionalProperties":false},"retentionPolicy":{"type":"array","items":{"title":"Space Retention Policy","type":"object","properties":{"spaceAdminCanEdit":{"type":"boolean"},"lastModifiedDate":{"type":"string"},"pageRetentionRule":{"$ref":"#/definitions/retention-rule"},"attachmentRetentionRule":{"$ref":"#/definitions/retention-rule"},"trashRetentionRule":{"title":"Trash Retention Rule","type":"object","properties":{"maxDeletedAge":{"type":"integer"},"deletedAgeUnit":{"type":"string","enum":["DAYS","MONTHS","YEARS"]},"keepAll":{"type":"boolean"}},"additionalProperties":false,"required":["keepAll"]},"lastModifiedBy":{"type":"string"}},"additionalProperties":false,"required":["spaceAdminCanEdit"]}}},"additionalProperties":false},"unknown-user":{"title":"Unknown User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"user":{"title":"User","type":"object","properties":{"profilePicture":{"$ref":"#/definitions/icon"},"displayName":{"type":"string"},"username":{"type":"string"}},"additionalProperties":false},"version":{"title":"Version","type":"object","properties":{"by":{"$ref":"#/definitions/person"},"when":{"type":"string"},"message":{"type":"string"},"number":{"type":"integer"},"minorEdit":{"type":"boolean"},"hidden":{"type":"boolean"},"syncRev":{"type":"string"},"content":{"type":"array","items":{"$ref":"#/definitions/content"}}},"additionalProperties":false,"required":["number","minorEdit","hidden"]},"web-resource-dependencies":{"title":"Web Resource Dependencies","type":"object","properties":{"keys":{"type":"array","items":{"type":"string"}},"contexts":{"type":"array","items":{"type":"string"}},"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"superbatch":{"type":"array","items":{"title":"Super Batch Web Resources","type":"object","properties":{"uris":{"type":"object","patternProperties":{".+":{"type":"array","items":{"type":"string","format":"uri"}}},"additionalProperties":false},"tags":{"type":"object","patternProperties":{".+":{"$ref":"#/definitions/html-string"}},"additionalProperties":false},"metatags":{"type":"array","items":{"$ref":"#/definitions/html-string"}}},"additionalProperties":false}}},"additionalProperties":false}},"additionalProperties":false}
```

  - **Status 404**

  - Returned if space exists for the specified space key or if the calling user does not have permission to perform the operation

#### [Remove space watch](#api/user/watch-removeSpaceWatch)`DELETE /rest/api/user/watch/space/{spaceKey}`

Delete an existing watcher for the given user and space key. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `DELETE http://example.com/confluence/rest/api/user/watch/space/SPACEKEY` 

- `DELETE http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?username=jblogs` 

- `DELETE http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?key=ff8080815a58e24c015a58e263710000` 

- `DELETE http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?contentType=blogpost` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to delete the watcher for |
| `username` | *string* | username of the user to delete the watcher for |
| `contentType` | *string* | the optional content type to delete the watcher for |

##### Responses

  - **Status 204** - *application/json*

  - Returned if the watcher was successfully deleted 

  - **Status 404**

  - Returned if no content exists for the specified space key or the calling user does not have permission to perform the operation

#### [Is watching space](#api/user/watch-isWatchingSpace)`GET /rest/api/user/watch/space/{spaceKey}`

Get information about whether a user is watching a specified space. 

User is optional. If not specified, currently logged-in user will be used. Otherwise, it can be specified by either user key or username. When a user is specified and is different from the logged-in user, the logged-in user needs to be a Confluence administrator.

Example request URI(s):

- `http://example.com/confluence/rest/api/user/watch/space/SPACEKEY` 

- `http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?username=jblogs` 

- `http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?key=ff8080815a58e24c015a58e263710000` 

- `http://example.com/confluence/rest/api/user/watch/space/SPACEKEY?contentType=blostpost` 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `key` | *string* | userkey of the user to check for watching state |
| `username` | *string* | username of the user to check for watching state |
| `contentType` | *string* | an optional content type to check for watching state |

##### Responses

  - **Status 200** - *application/json*

  - Returns a JSON representation containing the watching state 

  - **Status 404**

  - Returned if no space exists for the specified space key or the calling user does not have permission to perform the operation

### [api/webhooks](#api/webhooks)

#### [Create webhook](#api/webhooks-createWebhook)`POST /rest/api/webhooks`

Create a webhook via the URL. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-webhook#","title":"Rest Webhook","type":"object"}
```

##### Responses

  - **Status 201** - *application/json*

  - A created webhook. 

  - **Status 400** - *application/json*

  - The webhook parameters were invalid or not supplied.

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to create webhooks

#### [Find webhooks](#api/webhooks-findWebhooks)`GET /rest/api/webhooks`

Find webhooks. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `event` | *string* | list of \{@link com.atlassian.webhooks.WebhookEvent\} ids to filter for |
| `statistics` | *boolean* Default: `false` | \{@code true\} if statistics should be provided for all found webhooks |
| `start` | *int* |  |
| `limit` | *int* Default: `100` |  |

##### Responses

  - **Status 200** - *application/json*

  - A list of webhooks. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to find webhooks.

#### [Get webhook](#api/webhooks-getWebhook)`GET /rest/api/webhooks/{webhookId}`

Get a webhook by id. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `statistics` | *boolean* Default: `false` |  |

##### Responses

  - **Status 200** - *application/json*

  - a webhook. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to get a webhook.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Update webhook](#api/webhooks-updateWebhook)`PUT /rest/api/webhooks/{webhookId}`

Update an existing webhook. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### Schema

```
{"id":"https://docs.atlassian.com/jira/REST/schema/rest-webhook#","title":"Rest Webhook","type":"object"}
```

##### Responses

  - **Status 200** - *application/json*

  - a webhook. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to update a webhook.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Delete webhook](#api/webhooks-deleteWebhook)`DELETE /rest/api/webhooks/{webhookId}`

Delete a webhook via the URL. 

The authenticated user must be an administrator to call this resource. 

##### Responses

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to delete webhooks.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Get latest invocation](#api/webhooks-getLatestInvocation)`GET /rest/api/webhooks/{webhookId}/latest`

Get the latest invocations for a specific webhook. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `event` | *string* | the string id of a specific event to retrieve the last invocation for. |
| `outcome` | *string* | the outcome to filter for. Can be SUCCESS, FAILURE, ERROR. None specified means that the all will be considered  |

##### Responses

  - **Status 200** - *application/json*

  - a webhook invocation dataset. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to get webhook.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Get statistics](#api/webhooks-getStatistics)`GET /rest/api/webhooks/{webhookId}/statistics`

Get the statistics for a specific webhook. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `event` | *string* | the string id of a specific event to retrieve the last invocation for. May be empty, in which case all events are considered  |

##### Responses

  - **Status 200** - *application/json*

  - a webhook invocation dataset. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to get webhook.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Get statistics summary](#api/webhooks-getStatisticsSummary)`GET /rest/api/webhooks/{webhookId}/statistics/summary`

Get the statistics summary for a specific webhook. 

The authenticated user must be an administrator to call this resource. 

##### Responses

  - **Status 200** - *application/json*

  - a webhook invocation dataset. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to get webhook.

  - **Status 404** - *application/json*

  - The webhook does not exist.

#### [Test webhook](#api/webhooks-testWebhook)`POST /rest/api/webhooks/test`

Test connectivity to a specific endpoint. 

The authenticated user must be an administrator to call this resource. 

##### Request

###### query parameters

| **parameter** | **type** | **description** |
| :-: | :-: | :-: |
| `url` | *string* | the url in which to connect to |

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"request":{"url":"http://example.com/callback","body":null,"method":"POST","headers":{"Accept":"*/*"}},"response":{"status":200,"headers":{"Content-Type":"application/json"},"body":"{\n  \"accepted\": true\n}"}}
```

  - a webhook. 

  - **Status 401** - *application/json*

  - The currently authenticated user has insufficient permissions to test a connection. 

  - **Status 404** - *application/json*

  - The repository does not exist.

### [zdu](#zdu)

Provides Zero Downtime Upgrade capabilities.

#### [Approve upgrade](#zdu-approveUpgrade)`POST /rest/zdu/approve`

Finalizes the ZDU upgrade and runs specific tasks such as cleanup scripts.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"upgradeModeEnabled":false,"hasFinalizationTasks":true,"state":"STABLE","nodes":[{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}],"links":[{"rel":"self","href":"/"},{"rel":"start","href":"/start"}]}
```

  - Returns full JSON representation of the cluster 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster#","title":"Cluster","type":"object","properties":{"upgradeModeEnabled":{"type":"boolean"},"hasFinalizationTasks":{"type":"boolean"},"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"originalVersion":{"type":"string"},"nodes":{"type":"array","items":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"additionalProperties":false,"required":["local"]}},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"definitions":{"link":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}},"additionalProperties":false,"required":["upgradeModeEnabled","hasFinalizationTasks"]}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

  - **Status 409**

  - Returned if the cluster is not in a valid state

#### [Cancel upgrade](#zdu-cancelUpgrade)`POST /rest/zdu/cancel`

Prohibits the Upgrading of individual Nodes within the Cluster. All Nodes need to be on the same version before performing this request. 

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"upgradeModeEnabled":false,"hasFinalizationTasks":true,"state":"STABLE","nodes":[{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}],"links":[{"rel":"self","href":"/"},{"rel":"start","href":"/start"}]}
```

  - Returns full JSON representation of the cluster 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster#","title":"Cluster","type":"object","properties":{"upgradeModeEnabled":{"type":"boolean"},"hasFinalizationTasks":{"type":"boolean"},"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"originalVersion":{"type":"string"},"nodes":{"type":"array","items":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"additionalProperties":false,"required":["local"]}},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"definitions":{"link":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}},"additionalProperties":false,"required":["upgradeModeEnabled","hasFinalizationTasks"]}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

  - **Status 409**

  - Returned if the cluster is not in a valid state

#### [Get cluster](#zdu-getCluster)`GET /rest/zdu/cluster`

Gets an overview of a Cluster including its current state and composition of Nodes.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"upgradeModeEnabled":false,"hasFinalizationTasks":true,"state":"STABLE","nodes":[{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}],"links":[{"rel":"self","href":"/"},{"rel":"start","href":"/start"}]}
```

  - Returns full JSON representation of the cluster 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster#","title":"Cluster","type":"object","properties":{"upgradeModeEnabled":{"type":"boolean"},"hasFinalizationTasks":{"type":"boolean"},"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"originalVersion":{"type":"string"},"nodes":{"type":"array","items":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"additionalProperties":false,"required":["local"]}},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"definitions":{"link":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}},"additionalProperties":false,"required":["upgradeModeEnabled","hasFinalizationTasks"]}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

#### [Get node by id](#zdu-getNodeById)`GET /rest/zdu/nodes/{nodeId}`

Gets the requested Node's information.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}
```

  - Returns full JSON representation of the cluster node 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/node-info#","title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}}},"additionalProperties":false,"required":["local"]}
```

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

  - **Status 404**

  - Returned if there is no content with the given id

#### [Retry finalization](#zdu-retryFinalization)`POST /rest/zdu/retryUpgrade`

Reruns any incomplete finalization Upgrade Tasks.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"upgradeModeEnabled":true,"hasFinalizationTasks":true,"state":"READY_TO_UPGRADE","originalVersion":"1.0","nodes":[{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}],"links":[{"rel":"self","href":"/"},{"rel":"cancel","href":"/cancel"}]}
```

  - Returns full JSON representation of the cluster 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster#","title":"Cluster","type":"object","properties":{"upgradeModeEnabled":{"type":"boolean"},"hasFinalizationTasks":{"type":"boolean"},"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"originalVersion":{"type":"string"},"nodes":{"type":"array","items":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"additionalProperties":false,"required":["local"]}},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"definitions":{"link":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}},"additionalProperties":false,"required":["upgradeModeEnabled","hasFinalizationTasks"]}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

  - **Status 409**

  - Returned if the cluster is not in a valid state

#### [Start upgrade](#zdu-startUpgrade)`POST /rest/zdu/start`

Enables Upgrading of individual Nodes within the Cluster, allowing a heterogeneous Cluster formation.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"upgradeModeEnabled":true,"hasFinalizationTasks":true,"state":"READY_TO_UPGRADE","originalVersion":"1.0","nodes":[{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}],"links":[{"rel":"self","href":"/"},{"rel":"cancel","href":"/cancel"}]}
```

  - Returns full JSON representation of the cluster 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster#","title":"Cluster","type":"object","properties":{"upgradeModeEnabled":{"type":"boolean"},"hasFinalizationTasks":{"type":"boolean"},"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"originalVersion":{"type":"string"},"nodes":{"type":"array","items":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"additionalProperties":false,"required":["local"]}},"links":{"type":"array","items":{"$ref":"#/definitions/link"}}},"definitions":{"link":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}},"additionalProperties":false,"required":["upgradeModeEnabled","hasFinalizationTasks"]}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

  - **Status 409**

  - Returned if the cluster is not in a valid state

#### [Get state](#zdu-getState)`GET /rest/zdu/state`

Gets the State of the Cluster and the responding Node's information.

##### Responses

  - **Status 200** - *application/json*

###### Example

  - ```
{"state":"STABLE","buildInfo":{"id":"7a52c813","name":"node-1","ipAddress":"10.0.0.1","state":"ACTIVE","tasksTotal":5,"activeUserCount":12,"buildNumber":"1234","version":"1.0","local":true,"links":[{"rel":"self","href":"/nodes/7a52c813"}]}}
```

  - Returns full JSON representation of the cluster state 

###### Schema

  - ```
{"id":"https://docs.atlassian.com/jira/REST/schema/cluster-state-response#","title":"Cluster State Response","type":"object","properties":{"state":{"type":"string","enum":["STABLE","READY_TO_UPGRADE","MIXED","READY_TO_RUN_UPGRADE_TASKS","RUNNING_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"buildInfo":{"title":"Node Info","type":"object","properties":{"id":{"type":"string"},"name":{"type":"string"},"ipAddress":{"type":"string"},"state":{"type":"string","enum":["STARTING","ACTIVE","DRAINING","TERMINATING","OFFLINE","ERROR","RUNNING_FINALIZE_UPGRADE_TASKS","UPGRADE_TASKS_FAILED"]},"tasksTotal":{"type":"integer"},"activeUserCount":{"type":"integer"},"buildNumber":{"type":"string"},"version":{"type":"string"},"local":{"type":"boolean"},"finalization":{"title":"Node Finalization Info","type":"object","properties":{"errors":{"type":"array","items":{"title":"Upgrade Task Error","type":"object","properties":{"clusterUpgradeTask":{"type":"boolean"},"errors":{"type":"array","items":{"type":"string"}},"exceptionMessage":{"type":"string"}},"additionalProperties":false,"required":["clusterUpgradeTask"]}},"lastRequested":{"type":"string"}},"additionalProperties":false},"links":{"type":"array","items":{"title":"Link","type":"object","properties":{"rel":{"type":"string"},"href":{"type":"string","format":"uri"}},"additionalProperties":false}}},"additionalProperties":false,"required":["local"]}},"additionalProperties":false}
```

  - **Status 401**

  - Returned if user is not authenticated

  - **Status 403**

  - Returned if the calling user does not have permission to view the content

[Atlassian](http://www.atlassian.com/)
