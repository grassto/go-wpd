package gowpd

import (
	"log"
)

func Example_deleteFromDevice() {
	Initialize()

	mng, err := CreatePortableDeviceManager()
	if err != nil {
		panic(err)
	}

	devices, err := mng.GetDevices()
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

	// objectID which will be deleted from the device.
	targetObjectID := "F:\\test.txt"

	for _, id := range devices {
		portableDevice, err := CreatePortableDevice()
		if err != nil {
			panic(err)
		}

		err = portableDevice.Open(id, clientInfo)
		if err != nil {
			panic(err)
		}

		content, err := portableDevice.Content()
		if err != nil {
			panic(err)
		}

		pObjectsToDelete, err := CreatePortableDevicePropVariantCollection()

		pv := new(PropVariant)
		pv.Init()
		pv.Set(targetObjectID)
		err = pObjectsToDelete.Add(pv)
		if err != nil {
			panic(err)
		}
		results, err := content.Delete(PORTABLE_DEVICE_DELETE_NO_RECURSION, pObjectsToDelete)
		if err != nil {
			count, err := results.GetCount()
			if err != nil {
				panic(err)
			}
			log.Printf("Count: %d\n", count)
			result, err := results.GetAt(0)
			if err != nil {
				panic(err)
			}
			log.Printf("Type: %d\n", result.GetType())
			if result.GetType() == VT_ERROR {
				log.Printf("error: %#x\n", result.GetError())
			}

			panic(err)
		}

		pv.Clear()
		FreeDeviceID(id)
	}

	mng.Release()
	Uninitialize()

	// Output:
}
