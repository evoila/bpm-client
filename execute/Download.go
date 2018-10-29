package execute

import (
	. "github.com/evoila/BPM-Client/model"
	"github.com/evoila/BPM-Client/rest"
	"github.com/evoila/BPM-Client/s3"
)

func Download(requestBody PackageRequestBody) {

	var permission = rest.GetDownloadPermission(requestBody)

	s3.DownloadFile(requestBody.Name, permission)
}
