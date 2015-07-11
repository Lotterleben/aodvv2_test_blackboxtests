package main

import (
    "aodvv2_test_management"
    "fmt"
)

/* route states */
const
(
    ROUTE_STATE_ACTIVE = iota
    ROUTE_STATE_IDLE = iota
    ROUTE_STATE_INVALID = iota
    ROUTE_STATE_TIMED = iota
)

// TODO: unify those JSONs!!!! (seqnum vs orig_addr_seqnum usw)


type Json_template_sent_rreq struct {
    Log_type  string          `json: "log_type"`
    Log_data  struct {
        Orig_addr   string    `json: "orig_addr"`
        Targ_addr   string    `json: "targ_addr"`
        Orig_seqnum int       `json: "orig_seqnum"`
        Metric      int       `json: "metric"`
    } `json: "log_data"`
}

type Json_template_received_rreq struct {
    Log_type  string          `json: "log_type"`
    Log_data  struct {
        Last_hop    string    `json: "last_hop"`
        Orig_addr   string    `json: "orig_addr"` // TODO this is where I left off
        Targ_addr   string    `json: "targ_addr"`
        Orig_seqnum int       `json: "orig_seqnum"`
        Metric      int       `json: "metric"`
    } `json: "log_data"`
}

type Json_template_sent_rrep struct {

}

const test_string = "xoxotesttest"
const json_template_sent_rreq = "{\"log_type\": \"sent_rreq\", \"log_data\": {\"orig_addr\": \"%s\", \"targ_addr\": \"%s\", \"orig_seqnum\": %d, \"metric\": %d}}"
const json_template_received_rreq = "{\"log_type\": \"received_rreq\", \"log_data\":{\"last_hop\": \"%s\", \"orig_addr\": \"%s\", \"targ_addr\": \"%s\", \"orig_seqnum\": %d, \"metric\": %d}}"
const json_template_sent_rrep = "{\"log_type\": \"sent_rrep\", \"log_data\": {\"next_hop\": \"%s\",\"orig_addr\": \"%s\", \"orig_seqnum\": %d, \"targ_addr\": \"%s\"}}"
const json_template_received_rrep = "{\"log_type\": \"received_rrep\", \"log_data\":{\"last_hop\": \"%s\", \"orig_addr\": \"%s\", \"orig_seqnum\": %d, \"targ_addr\": \"%s\", \"targ_addr_seqnum\": %d}}"
const json_template_added_rt_entry = "{\"log_type\": \"added_rt_entry\", \"log_data\": {\"addr\": \"%s\", \"next_hop\": \"%s\", \"seqnum\": %d, \"metric\": %d, \"state\": %d}}"


func test_route_creation_0_to_3() {
    test_name := "route_creation_0_to_3"

    riot_line := aodvv2_test_management.Create_clean_setup(test_name)

    fmt.Println("Starting test ", test_name)

    beginning := riot_line[0]
    end := riot_line[len(riot_line)-1]

    beginning.Channels.Send(fmt.Sprintf("send %s %s\n", end.Ip, test_string))

    /* Discover route...  */
    expected_json := fmt.Sprintf(json_template_sent_rreq, beginning.Ip, end.Ip, 1, 0)
    beginning.Channels.Expect_JSON(expected_json)

    xoxo := json.Marshal(&Json_template_received_rreq{
        Log_type: "received_rreq",
        Log_data: struct {
            Last_hop:    beginning.Ip,
            Orig_addr:   beginning.Ip,
            Targ_addr:   end.Ip,
            Orig_seqnum: 1,
            Metric:      0}})
    //fmt.Println(string(xoxo))
    riot_line[1].Channels.Expect_JSON(string(xoxo))

    expected_json = fmt.Sprintf(json_template_received_rreq, beginning.Ip, beginning.Ip, end.Ip, 1, 0)
    riot_line[1].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.Ip, beginning.Ip, 1, 1, ROUTE_STATE_ACTIVE)
    riot_line[1].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_sent_rreq, beginning.Ip, end.Ip, 1, 1)
    riot_line[1].Channels.Expect_JSON(expected_json)

    expected_json = fmt.Sprintf(json_template_received_rreq, riot_line[1].Ip, beginning.Ip, end.Ip, 1, 1)
    riot_line[2].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.Ip, riot_line[1].Ip, 1, 2, ROUTE_STATE_ACTIVE)
    riot_line[2].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_sent_rreq, beginning.Ip, end.Ip, 1, 2)
    riot_line[2].Channels.Expect_JSON(expected_json)

    expected_json = fmt.Sprintf(json_template_received_rreq, riot_line[2].Ip, beginning.Ip, end.Ip, 1, 2)
    end.Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, beginning.Ip, riot_line[2].Ip, 1, 3, ROUTE_STATE_ACTIVE)
    end.Channels.Expect_JSON(expected_json)

    /* And send a RREP back */

    expected_json = fmt.Sprintf(json_template_sent_rrep, riot_line[2].Ip, beginning.Ip, 1, end.Ip)
    end.Channels.Expect_JSON(expected_json)

    /* TODO FIXME targ_addr_seqnum is 2 because apparently
     * weird RREQs are sent out before the experiment, screwing up the targaddr seqnum
     * and I haven't figured out why yet. */
    expected_json= fmt.Sprintf(json_template_received_rrep, end.Ip, beginning.Ip, 1, end.Ip, 2)
    riot_line[2].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, end.Ip, end.Ip, 2, 1, ROUTE_STATE_ACTIVE)
    riot_line[2].Channels.Expect_JSON(expected_json)
    expected_json= fmt.Sprintf(json_template_sent_rrep, riot_line[1].Ip, beginning.Ip, 1, end.Ip)
    riot_line[2].Channels.Expect_JSON(expected_json)

    expected_json= fmt.Sprintf(json_template_received_rrep, riot_line[2].Ip, beginning.Ip, 1, end.Ip, 2)
    riot_line[1].Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, end.Ip, riot_line[2].Ip, 1, 2, ROUTE_STATE_ACTIVE)
    riot_line[1].Channels.Expect_JSON(expected_json)
    expected_json= fmt.Sprintf(json_template_sent_rrep, beginning.Ip, beginning.Ip, 1, end.Ip)
    riot_line[1].Channels.Expect_JSON(expected_json)

    expected_json= fmt.Sprintf(json_template_received_rrep, riot_line[1].Ip, beginning.Ip, 1, end.Ip, 2)
    beginning.Channels.Expect_JSON(expected_json)
    expected_json = fmt.Sprintf(json_template_added_rt_entry, end.Ip, riot_line[1].Ip, 1, 3, ROUTE_STATE_ACTIVE)
    beginning.Channels.Expect_JSON(expected_json)

    //TODO: defer dump Channels
    fmt.Println("\nDone.")
}

func test_route_creation_3_to_0() {
    /*
    test_name := "route_creation_0_to_3"

    riot_line := aodvv2_test_management.Create_clean_setup(test_name)

    fmt.Println("Starting test ", test_name)

    beginning := riot_line[3]
    end := riot_line[0]

    beginning.Channels.Send(fmt.Sprintf("send %s %s\n", end.Ip, test_string))


    fmt.Println("\nDone.")
    */
}

func start_experiments() {
    test_route_creation_0_to_3()
    test_route_creation_3_to_0()
}

func main() {
    //TODO: build fresh RIOT image
    start_experiments()
}