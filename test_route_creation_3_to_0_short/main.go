package main

import (
    mgmt "aodvv2_test_management"
    "fmt"
)

func check(e error) {
    if e != nil {
        fmt.Println("OMG EVERYBODY PANIC")
        panic(e)
    }
}


/* Basically the same as route_creation_3_to_0: send teststring from 3 to 0 and
 * see if a route is established, but this time without any intermediate checks
 * (since they seem to avoid breaking behavior that has been seen in the wild)
 */
func main() {
    test_name := "route_creation_3_to_0_short"
    riot_line := mgmt.Create_clean_setup(test_name)

    fmt.Println("Starting test ", test_name)

    /* note: the seqnums are *sometimes* != 1 because apparently
     * weird RREQs are sent out before the experiment, screwing up the node's seqnum
     * and I haven't figured out why yet. So until that mystery is solved
     * they are set to aodvv2_test_management.WILDCARD */
    orig_seqnum := "\""+ mgmt.WILDCARD+ "\""
    targ_seqnum := "\""+ mgmt.WILDCARD+ "\""
    beginning := riot_line[3]
    end := riot_line[0]

    /* port number (1234) doesn't matter since no one is listening at that port
       because we don't care about anyone receiving the actual content */
    beginning.Channels.Send(fmt.Sprintf("udp send %s 1234 %s\n", end.Ip, mgmt.Test_string))
    /* Now a route should be discovered.....*/

    /* check node 3 (aka the beginning) */
    beginning.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rrep, map[string]string{
        "Last_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
        "Targ_seqnum": targ_seqnum,
    }))

    fmt.Println("\nDone.")

}