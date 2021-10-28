package gowpd

import (
	"log"
)

func RecursiveEnumerate(parentObjectID string, content *IPortableDeviceContent) {
	enum, err := content.EnumObjects(parentObjectID)
	if err != nil {
		panic(err)
	}

	objectIDs := make([]string, 0)
	for {
		tmp, err := enum.Next(10)
		if err != nil {
			panic(err)
		}
		if len(tmp) == 0 {
			break
		}
		objectIDs = append(objectIDs, tmp...)
	}

	for _, objectID := range objectIDs {
		log.Println(objectID)
	}

	for _, objectID := range objectIDs {
		RecursiveEnumerate(objectID, content)
	}
}

func Example_contentEnumerate() {
	Initialize()

	mng, err := CreatePortableDeviceManager()
	if err != nil {
		panic(err)
	}

	pClientInfo, err := CreatePortableDeviceValues()
	if err != nil {
		panic(err)
	}
	pClientInfo.SetStringValue(WPD_CLIENT_NAME, "libgowpd")
	pClientInfo.SetUnsignedIntegerValue(WPD_CLIENT_MAJOR_VERSION, 1)
	pClientInfo.SetUnsignedIntegerValue(WPD_CLIENT_MINOR_VERSION, 0)
	pClientInfo.SetUnsignedIntegerValue(WPD_CLIENT_REVISION, 2)

	deviceIDs, err := mng.GetDevices()
	if err != nil {
		panic(err)
	}

	for _, deviceID := range deviceIDs {
		device, err := CreatePortableDevice()
		if err != nil {
			panic(err)
		}

		err = device.Open(deviceID, pClientInfo)
		if err != nil {
			panic(err)
		}

		content, err := device.Content()
		if err != nil {
			panic(err)
		}

		RecursiveEnumerate(WPD_DEVICE_OBJECT_ID, content)

		FreeDeviceID(deviceID)
	}

	Uninitialize()
}
