#include "stdafx.h"
#include "libgowpd.h"
#include <string>
#include <iostream>

std::string WCHAR2String(wchar_t* pwszSrc);
void Wchar_tToString(std::string& szDst, wchar_t* wchar);
void StringToWstring(std::wstring& szDst, std::string str);

int main()
{
	HRESULT hr;

	hr = CoInitializeEx(NULL, COINIT_MULTITHREADED);
	if (FAILED(hr)) {
		printf("init err %x\n", hr);
	}
	IPortableDeviceManager *pPortableDeviceManager;
	hr = CoCreateInstance(
		CLSID_PortableDeviceManager,
		NULL,
		CLSCTX_INPROC_SERVER,
		IID_IPortableDeviceManager,
		(LPVOID*)&pPortableDeviceManager);
	if (SUCCEEDED(hr)) {
		PnPDeviceID *pPnPDeviceIDs;
		DWORD cPnPDeviceIDs;
		hr = portableDeviceManager_GetDevices(pPortableDeviceManager, &pPnPDeviceIDs, &cPnPDeviceIDs);
		if (SUCCEEDED(hr)) {
			printf("%d devices has been found.\n", cPnPDeviceIDs);

			for (int i = 0; i < cPnPDeviceIDs; i++) {
				PnPDeviceID id = pPnPDeviceIDs[i];

				std::string szDst = WCHAR2String(id);
				std::cout << szDst << std::endl;

				PWSTR friendlyName;
				PWSTR manufacturer;
				PWSTR description;
				DWORD cFriendlyName = 0;
				DWORD cManufacturer = 0;
				DWORD cDescription = 0;

				hr = portableDeviceManager_GetDeviceFriendlyName(pPortableDeviceManager, id, &friendlyName, &cFriendlyName);
				if (SUCCEEDED(hr)) {
					printf("%ws\n", friendlyName);
					free(friendlyName);
					friendlyName = NULL;
				} else {
					printf("friendlyname wtf 0x%x\n", hr);
				}
				hr = portableDeviceManager_GetDeviceManufacturer(pPortableDeviceManager, id, &manufacturer, &cManufacturer);
				if (SUCCEEDED(hr)) {
					printf("%ws\n", manufacturer);
					free(manufacturer);
					manufacturer = NULL;
				} else {
					printf("manufacturer WTF 0x%x\n", hr);
				}
				hr = portableDeviceManager_GetDeviceDescription(pPortableDeviceManager, id, &description, &cDescription);
				if (SUCCEEDED(hr)) {
					printf("%ws\n", description);
					free(description);
					description = NULL;
				} else {
					printf("description WTF 0x%x\n", hr);
				}

				CoTaskMemFree(id);
			}

			free(pPnPDeviceIDs);
		}
		else {
			printf("get devices WTF %x\n", hr);
		}
	}
	else {
		printf("create manager WTF %x\n", hr);
	}

	PROPVARIANT prop = {0};
	PropVariantInit(&prop);

	prop.vt = VT_LPWSTR;
	prop.wReserved1 = 2;
	prop.wReserved2 = 3;
	prop.wReserved3 = 4;
	prop.pwszVal = (LPWSTR)CoTaskMemAlloc(sizeof(WCHAR) * 5);
	char* testStr = "TEST";
	for (int i = 0; i < 4; i++) {
		prop.pwszVal[i] = testStr[i];
	}
	prop.pwszVal[4] = 0;

	printf("prop.vt:      0x%x\n", prop.vt);
	printf("prop.pwszVal: 0x%016p\n", prop.pwszVal);
	byte* bs = (byte*) &prop;
	for (int i = 0; i < sizeof(PROPVARIANT); i++) {
		printf("0x%x\n", bs[i]);
	}
	printf("LPWSTR size:  0x%x\n", sizeof(LPWSTR));

	PropVariantClear(&prop);

    return 0;
}

// WCHAR 转换为 std::string
std::string WCHAR2String(wchar_t* pwszSrc)
{
	int nLen = WideCharToMultiByte(CP_ACP, 0, pwszSrc, -1, NULL, 0, NULL, NULL);
	if (nLen <= 0)
		return std::string("");

	char* pszDst = new char[nLen];
	if (NULL == pszDst)
		return std::string("");

	WideCharToMultiByte(CP_ACP, 0, pwszSrc, -1, pszDst, nLen, NULL, NULL);
	pszDst[nLen - 1] = 0;

	std::string strTmp(pszDst);
	delete[] pszDst;

	return strTmp;
}

void Wchar_tToString(std::string& szDst, wchar_t* wchar)
{
	wchar_t* wText = wchar;
	DWORD dwNum = WideCharToMultiByte(CP_OEMCP, NULL, wText, -1, NULL, 0, NULL, FALSE);// WideCharToMultiByte的运用
	char* psText; // psText为char*的临时数组，作为赋值给std::string的中间变量
	psText = new char[dwNum];
	WideCharToMultiByte(CP_OEMCP, NULL, wText, -1, psText, dwNum, NULL, FALSE);// WideCharToMultiByte的再次运用
	szDst = psText;// std::string赋值
	delete[]psText;// psText的清除
}

void StringToWstring(std::wstring& szDst, std::string str)
{
	std::string temp = str;
	int len = MultiByteToWideChar(CP_ACP, 0, (LPCSTR)temp.c_str(), -1, NULL, 0);
	wchar_t* wszUtf8 = new wchar_t[len + 1];
	memset(wszUtf8, 0, len * 2 + 2);
	MultiByteToWideChar(CP_ACP, 0, (LPCSTR)temp.c_str(), -1, (LPWSTR)wszUtf8, len);
	szDst = wszUtf8;
	std::wstring r = wszUtf8;
	delete[] wszUtf8;
}