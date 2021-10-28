package gowpd

import (
	"log"
)

func Example_transferToPC() {
	Initialize()

	mng, err := CreatePortableDeviceManager()
	if err != nil {
		panic(err)
	}

	deviceIDs, err := mng.GetDevices()
	if err != nil {
		panic(err)
	}

	clientInfo, err := CreatePortableDeviceValues()
	if err != nil {
		panic(err)
	}
	clientInfo.SetStringValue(WPD_CLIENT_NAME, "libgowpd")
	clientInfo.SetUnsignedIntegerValue(WPD_CLIENT_MAJOR_VERSION, 1)
	clientInfo.SetUnsignedIntegerValue(WPD_CLIENT_MINOR_VERSION, 0)
	clientInfo.SetUnsignedIntegerValue(WPD_CLIENT_REVISION, 2)

	// object ID which will be transferred to PC.
	targetObjectID := "F:\\test.txt"
	// location where file will be transferred into.
	targetDestination := "E:\\test.txt"

	for _, id := range deviceIDs {
		portableDevice, err := CreatePortableDevice()
		if err != nil {
			panic(err)
		}

		portableDevice.Open(id, clientInfo)

		content, err := portableDevice.Content()
		if err != nil {
			panic(err)
		}
		resources, err := content.Transfer()
		if err != nil {
			panic(err)
		}

		objectDataStream, optimalTransferSize, err := resources.GetStream(targetObjectID, WPD_RESOURCE_DEFAULT, STGM_READ)
		if err != nil {
			panic(err)
		}

		pFinalFileStream, err := SHCreateStreamOnFile(targetDestination, STGM_CREATE|STGM_WRITE)
		if err != nil {
			panic(err)
		}

		totalBytesWritten, err := StreamCopy(pFinalFileStream, objectDataStream, optimalTransferSize)
		if err != nil {
			panic(err)
		}

		err = pFinalFileStream.Commit(0)
		if err != nil {
			panic(err)
		}

		log.Printf("Total bytes written: %d\n", totalBytesWritten)

		FreeDeviceID(id)
		portableDevice.Release()
	}

	mng.Release()
	Uninitialize()

	// Output:
}
