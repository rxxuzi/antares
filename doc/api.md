# API Documentation

This document describes the usage of the Antares API endpoints.

## Base URL

All API requests should be sent to: `/api`

## Request Format

All requests must be POST requests with a JSON payload. The general structure of a request is as follows:

- `file`: Boolean indicating whether the operation is on a file (true) or a folder (false)
- `type`: String indicating the type of operation
- `path`: String representing the path of the file or folder
- `src`: String representing the source path (for move and copy operations)
- `dst`: String representing the destination path (for move, copy, and rename operations)
- `option`: Object containing additional options (if any)

## Response Format

All responses are returned in JSON format with the following structure `success`: the operation succeeded:

- `success`: Boolean indicating whether the operation was successful
- `message`: String providing additional information about the result
- `error_code`: String indicating the type of error (if any)
- `data`: Object containing any additional data (if applicable)

## Available Operations

### 1. Delete

Deletes a file or folder.

**Request:**

```json
{
  "file": true,
  "type": "delete",
  "path": "/path/to/file.txt"
}
```

**Response:**

```json
{
  "success": true,
  "message": "File deleted successfully: /path/to/file.txt"
}
```

### 2. Move

Moves a file or folder to a new location.

**Request:**

```json
{
  "file": true,
  "type": "move",
  "path": "/path/to/source/file.txt",
  "dst": "/path/to/destination/file.txt"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Successfully moved from /path/to/source/file.txt to /path/to/destination/file.txt"
}
```

### 3. Copy

Creates a copy of a file.

**Request:**

```json
{
  "file": true,
  "type": "copy",
  "path": "/path/to/source/file.txt"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Successfully copied from /path/to/source/file.txt to /path/to/source/file-copy.txt"
}
```

### 4. Rename

Renames a file or folder.

**Request:**

```json
{
  "file": true,
  "type": "rename",
  "path": "/path/to/oldname.txt",
  "dst": "/path/to/newname.txt"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Successfully renamed from /path/to/oldname.txt to /path/to/newname.txt"
}
```

### 5. Create Directory (mkdir)

Creates a new directory.

**Request:**

```json
{
  "file": false,
  "type": "mkdir",
  "path": "/path/to/new/directory"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Directory created successfully: /path/to/new/directory"
}
```

## Error Codes

If an error occurs on the server side, the API will return one of the following error codes

- `INVALID_METHOD`: The HTTP method used is not allowed
- `INVALID_JSON`: The request body is not valid JSON
- `UNKNOWN_OPERATION`: The requested operation type is not recognized
- `MISSING_PATH`: The required path parameter is missing in the requested json
- `INVALID_PATH`: The specified path is invalid or outside the allowed directories
- `FILE_NOT_FOUND`: The specified file or folder does not exist
- `OPERATION_FAILED`: The requested operation failed
