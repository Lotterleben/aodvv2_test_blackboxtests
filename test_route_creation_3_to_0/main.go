package main

import (
    mgmt "aodvv2_test_management"
    "fmt"
    "strconv"
)

func check(e error) {
    if e != nil {
        fmt.Println("OMG EVERYBODY PANIC")
        panic(e)
    }
}

func main() {
    test_name := "route_creation_3_to_0"

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

    beginning.Channels.Send(fmt.Sprintf("send %s %s\n", end.Ip, mgmt.Test_string))

    /* Discover route at node 3...  */
    beginning.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "0",
    }))

    /* check node 2 */
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rreq, map[string]string{
        "Last_hop": beginning.Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "0",
    }))
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": beginning.Ip,
        "Seqnum": orig_seqnum,
        "Metric": "1",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "1",
    }))

    /* check node 1 */
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rreq, map[string]string{
        "Last_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "1",
    }))
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": riot_line[2].Ip,
        "Seqnum": orig_seqnum,
        "Metric": "2",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "2",
    }))

    /* check node 0 (aka the end) */
    end.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rreq, map[string]string{
        "Last_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": orig_seqnum,
        "Metric": "2",
    }))
    end.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": riot_line[1].Ip,
        "Seqnum": orig_seqnum,
        "Metric": "3",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))
    /* And send a RREP back */
    end.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rrep, map[string]string{
        "Next_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
    }))

    /* check node 1 */
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rrep, map[string]string{
        "Last_hop": end.Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
        "Targ_seqnum": "1",
    }))
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": end.Ip,
        "Seqnum": targ_seqnum,
        "Metric": "1",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))
    riot_line[1].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rrep, map[string]string{
        "Next_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
    }))

    /* check node 2 */
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rrep, map[string]string{
        "Last_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
        "Targ_seqnum": targ_seqnum,
    }))
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": riot_line[1].Ip,
        "Seqnum": targ_seqnum,
        "Metric": "2",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))
    riot_line[2].Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_sent_rrep, map[string]string{
        "Next_hop": beginning.Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
    }))

    /* check node 3 (aka the beginning) */
    beginning.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_received_rrep, map[string]string{
        "Last_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": orig_seqnum,
        "Targ_addr": end.Ip,
        "Targ_seqnum": targ_seqnum,
    }))
    beginning.Channels.Expect_JSON(mgmt.Make_JSON_str(mgmt.Tmpl_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": riot_line[2].Ip,
        "Seqnum": targ_seqnum,
        "Metric": "3",
        "State": strconv.Itoa(mgmt.ROUTE_STATE_ACTIVE),
    }))

    fmt.Println("\nDone.")
}