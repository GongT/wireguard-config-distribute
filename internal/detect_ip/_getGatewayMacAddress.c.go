// +build windows

package detect_ip

// #ifndef WIN32_LEAN_AND_MEAN
// #define WIN32_LEAN_AND_MEAN
// #endif

// #include <windows.h>
// #include <winsock2.h>
// #include <ws2ipdef.h>
// #include <iphlpapi.h>
// #include <stdio.h>
// #include <stdlib.h>

// #pragma comment(lib, "iphlpapi.lib")
// #pragma comment(lib, "ws2_32.lib")

// int TestFn()
// {
// 	int i;
// 	unsigned int j;
// 	unsigned long status = 0;

// 	PMIB_IPNET_TABLE2 pipTable = NULL;

// 	status = GetIpNetTable2(AF_INET, &pipTable);
// 	if (status != NO_ERROR) {
// 		printf("GetIpNetTable for IPv4 table returned error: %ld\n", status);
// 		exit(1);
// 	}
// 	printf("Number of IPv4 table entries: %d\n\n", pipTable->NumEntries);

// 	for (i = 0; (unsigned) i < pipTable->NumEntries; i++) {
// 		printf("IPv4 Address[%d]:\t %s\n", (int) i,
// 			inet_ntoa(pipTable->Table[i].Address.Ipv4.sin_addr));
// 		printf("Interface index[%d]:\t\t %lu\n", (int) i,
// 			pipTable->Table[i].InterfaceIndex);

// 		printf("Interface LUID NetLuidIndex[%d]:\t %lu\n",
// 			(int) i, pipTable->Table[i].InterfaceLuid.Info.NetLuidIndex);
// 		printf("Interface LUID IfType[%d]: ", (int) i);
// 		switch (pipTable->Table[i].InterfaceLuid.Info.IfType) {
// 		case IF_TYPE_OTHER:
// 			printf("Other\n");
// 			break;
// 		case IF_TYPE_ETHERNET_CSMACD:
// 			printf("Ethernet\n");
// 			break;
// 		case IF_TYPE_ISO88025_TOKENRING:
// 			printf("Token ring\n");
// 			break;
// 		case IF_TYPE_PPP:
// 			printf("PPP\n");
// 			break;
// 		case IF_TYPE_SOFTWARE_LOOPBACK:
// 			printf("Software loopback\n");
// 			break;
// 		case IF_TYPE_ATM:
// 			printf("ATM\n");
// 			break;
// 		case IF_TYPE_IEEE80211:
// 			printf("802.11 wireless\n");
// 			break;
// 		case IF_TYPE_TUNNEL:
// 			printf("Tunnel encapsulation\n");
// 			break;
// 		case IF_TYPE_IEEE1394:
// 			printf("IEEE 1394 (Firewire)\n");
// 			break;
// 		default:
// 			printf("Unknown: %d\n",
// 				pipTable->Table[i].InterfaceLuid.Info.IfType);
// 			break;
// 		}

// 		printf("Physical Address[%d]:\t ", (int) i);
// 		if (pipTable->Table[i].PhysicalAddressLength == 0)
// 			printf("\n");
// 		for (j = 0; j < pipTable->Table[i].PhysicalAddressLength; j++) {
// 			if (j == (pipTable->Table[i].PhysicalAddressLength - 1))
// 				printf("%.2X\n", (int) pipTable->Table[i].PhysicalAddress[j]);
// 			else
// 				printf("%.2X-", (int) pipTable->Table[i].PhysicalAddress[j]);
// 		}

// 		printf("Physical Address Length[%d]:\t %lu\n", (int) i,
// 			pipTable->Table[i].PhysicalAddressLength);

// 		printf("Neighbor State[%d]:\t ", (int) i);
// 		switch (pipTable->Table[i].State) {
// 		case NlnsUnreachable:
// 			printf("NlnsUnreachable\n");
// 			break;
// 		case NlnsIncomplete:
// 			printf("NlnsIncomplete\n");
// 			break;
// 		case NlnsProbe:
// 			printf("NlnsProbe\n");
// 			break;
// 		case NlnsDelay:
// 			printf("NlnsDelay\n");
// 			break;
// 		case NlnsStale:
// 			printf("NlnsStale\n");
// 			break;
// 		case NlnsReachable:
// 			printf("NlnsReachable\n");
// 			break;
// 		case NlnsPermanent:
// 			printf("NlnsPermanent\n");
// 			break;
// 		default:
// 			printf("Unknown: %d\n", pipTable->Table[i].State);
// 			break;
// 		}

// 		printf("Flags[%d]:\t\t %u\n", (int) i,
// 			(unsigned char) pipTable->Table[i].Flags);

// 		printf("ReachabilityTime[%d]:\t %lu\n\n", (int) i,
// 			pipTable->Table[i].ReachabilityTime);

// 	}
// 	FreeMibTable(pipTable);
// 	pipTable = NULL;

// 	exit(0);
// }
import "C"

func testFn() {
	C.TestFn()
}
