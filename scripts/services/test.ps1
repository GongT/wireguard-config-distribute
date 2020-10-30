
	$AccessRule = New-Object -TypeName System.Security.AccessControl.FileSystemAccessRule `
			-ArgumentList "NT AUTHORITY\NetworkService","FullControl","Allow"
