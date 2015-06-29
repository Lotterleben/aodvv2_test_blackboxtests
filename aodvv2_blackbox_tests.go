package main

import (
    "aodvv2_test_management"
    "fmt"
)

func test_route_creation_0_to_3() {
    /* route states */
    const
    (
        ROUTE_STATE_ACTIVE = iota
        ROUTE_STATE_IDLE = iota
        ROUTE_STATE_INVALID = iota
        ROUTE_STATE_TIMED = iota
    )

    const test_string = "xoxotesttest"
    const json_template_sent_rreq = "{\"log_type\": \"sent_rreq\", \"log_data\": {\"orig_addr\": \"%s\", \"targ_addr\": \"%s\", \"seqnum\": %d, \"metric\": %d}}"
    const json_template_received_rreq = "{\"log_type\": \"received_rreq\", \"log_data\":{\"last_hop\": \"%s\", \"orig_addr\": \"%s\", \"targ_addr\": \"%s\", \"orig_addr_seqnum\": %d, \"metric\": %d}}"
    const json_template_sent_rrep = "{\"log_type\": \"sent_rrep\", \"log_data\": {\"next_hop\": \"%s\",\"orig_addr\": \"%s\", \"orig_addr_seqnum\": %d, \"targ_addr\": \"%s\"}}"
    const json_template_received_rrep = "{\"log_type\": \"received_rrep\", \"log_data\":{\"last_hop\": \"%s\", \"orig_addr\": \"%s\", \"orig_addr_seqnum\": %d, \"targ_addr\": \"%s\"}}"
    const json_template_added_rt_entry = "{\"log_type\": \"added_rt_entry\", \"log_data\": {\"addr\": \"%s\", \"next_hop\": \"%s\", \"seqnum\": %d, \"metric\": %d, \"state\": %d}}"

    riot_line := aodvv2_test_management.Create_clean_setup("testtest")

    fmt.Println("Starting test...")

    beginning := riot_line[0]
    end := riot_line[len(riot_line)-1]

    beginning.Channels.Send(fmt.Sprintf("send %s %s\n", end.Ip, test_string))
    fmt.Print(".")

    /* Discover route...  */

    /*
    expected_json := fmt.Sprintf(json_template_sent_rreq, beginning.ip, end.ip, 1, 0)
    fmt.Print(".")
    beginning.channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json = fmt.Sprintf(json_template_received_rreq, beginning.ip, beginning.ip, end.ip, 1, 0)
    riot_line[1].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.ip, beginning.ip, 1, 1, ROUTE_STATE_ACTIVE)
    riot_line[1].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json = fmt.Sprintf(json_template_sent_rreq, beginning.ip, end.ip, 1, 1)
    riot_line[1].channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json = fmt.Sprintf(json_template_received_rreq, riot_line[1].ip, beginning.ip, end.ip, 1, 1)
    riot_line[2].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.ip, riot_line[1].ip, 1, 2, ROUTE_STATE_ACTIVE)
    riot_line[2].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json = fmt.Sprintf(json_template_sent_rreq, beginning.ip, end.ip, 1, 2)
    riot_line[2].channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json = fmt.Sprintf(json_template_received_rreq, riot_line[2].ip, beginning.ip, end.ip, 1, 2)
    end.channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.ip, riot_line[2].ip, 1, 3, ROUTE_STATE_ACTIVE)
    end.channels.expect_JSON(expected_json)
    fmt.Print(".")
    */
    /* And send a RREP back */
    /* NOTE: added_rt_entry isn't checked on the was back yet because apparently
     * weird RREQs are sent out before the experiment, screwing up the targaddr seqnum
     * and I haven't figured out why yet. TODO FIXME */

    /*
    expected_json = fmt.Sprintf(json_template_sent_rrep, riot_line[2].ip, beginning.ip, 1, end.ip)
    end.channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json= fmt.Sprintf(json_template_received_rrep, end.ip, beginning.ip, 1, end.ip)
    riot_line[2].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json= fmt.Sprintf(json_template_sent_rrep, riot_line[1].ip, beginning.ip, 1, end.ip)
    riot_line[2].channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json= fmt.Sprintf(json_template_received_rrep, riot_line[2].ip, beginning.ip, 1, end.ip)
    riot_line[1].channels.expect_JSON(expected_json)
    fmt.Print(".")
    expected_json= fmt.Sprintf(json_template_sent_rrep, beginning.ip, beginning.ip, 1, end.ip)
    riot_line[1].channels.expect_JSON(expected_json)
    fmt.Print(".")

    expected_json= fmt.Sprintf(json_template_received_rrep, riot_line[1].ip, beginning.ip, 1, end.ip)
    beginning.channels.expect_JSON(expected_json)
    fmt.Print(".")

    //TODO: defer dump channels
    fmt.Println("\nDone.")
    */
}

func start_experiments() {
    /* TODO: move this to dedicated test file */
    test_route_creation_0_to_3()
}

func main() {
    //TODO: build fresh RIOT image
    start_experiments()
}