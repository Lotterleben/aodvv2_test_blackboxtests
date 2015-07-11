package main

import (
    "aodvv2_test_management"
    "bytes"
    "fmt"
    "strconv"
    "text/template"
)

/* route states */
const
(
    ROUTE_STATE_ACTIVE = iota
    ROUTE_STATE_IDLE = iota
    ROUTE_STATE_INVALID = iota
    ROUTE_STATE_TIMED = iota
)

const test_string = "xoxotesttest"

const template_sent_rreq = "{\"log_type\": \"sent_rreq\", "+
                            "\"log_data\": {"+
                                    "\"orig_addr\": \"{{.Orig_addr}}\", "+
                                    "\"orig_seqnum\": {{.Orig_seqnum}}, "+
                                    "\"targ_addr\": \"{{.Targ_addr}}\", "+
                                    "\"metric\": {{.Metric}}}}"

const template_received_rreq = "{\"log_type\": \"received_rreq\", "+
                                "\"log_data\":{"+
                                    "\"last_hop\": \"{{.Last_hop}}\", "+
                                    "\"orig_addr\": \"{{.Orig_addr}}\", "+
                                    "\"orig_seqnum\": {{.Orig_seqnum}}, "+
                                    "\"targ_addr\": \"{{.Targ_addr}}\", "+
                                    "\"metric\": {{.Metric}}}}"

const template_added_rt_entry = "{\"log_type\": \"added_rt_entry\", "+
                                 "\"log_data\": {"+
                                    "\"addr\": \"{{.Addr}}\", "+
                                    "\"next_hop\": \"{{.Next_hop}}\", "+
                                    "\"seqnum\": {{.Seqnum}}, "+
                                    "\"metric\": {{.Metric}}, "+
                                    "\"state\": {{.State}}}}"

const template_sent_rrep = "{\"log_type\": \"sent_rrep\", "+
                            "\"log_data\": {"+
                                "\"next_hop\": \"{{.Next_hop}}\", "+
                                "\"orig_addr\": \"{{.Orig_addr}}\", "+
                                "\"orig_seqnum\": {{.Orig_seqnum}}, "+
                                "\"targ_addr\": \"{{.Targ_addr}}\"}}"

const template_received_rrep = "{\"log_type\": \"received_rrep\", "+
                            "\"log_data\":{"+
                                "\"last_hop\": \"{{.Last_hop}}\", "+
                                "\"orig_addr\": \"{{.Orig_addr}}\", "+
                                "\"orig_seqnum\": {{.Orig_seqnum}}, "+
                                "\"targ_addr\": \"{{.Targ_addr}}\", "+
                                "\"targ_seqnum\": {{.Targ_seqnum}}}}"

func check(e error) {
    if e != nil {
        fmt.Println("OMG EVERYBODY PANIC")
        panic(e)
    }
}

/* Create a JSON string from a given template (tmpl) and map containing the values
 * to be added to the template (data). */
func make_JSON_str(tmpl string, data map[string]string) string {
    strbuf := new(bytes.Buffer)
    t, _ := template.New("test").Parse(tmpl)
    // TODO: get writer to write to string, return string
    err := t.Execute(strbuf, data)
    check(err)
    return strbuf.String()
}


func test_route_creation_0_to_3() {
    test_name := "route_creation_0_to_3"

    riot_line := aodvv2_test_management.Create_clean_setup(test_name)

    fmt.Println("Starting test ", test_name)

    beginning := riot_line[0]
    end := riot_line[len(riot_line)-1]

    beginning.Channels.Send(fmt.Sprintf("send %s %s\n", end.Ip, test_string))

    /* Discover route at node 0...  */
    beginning.Channels.Expect_JSON(make_JSON_str(template_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "0",
    }))

    /* check node 1 */
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_received_rreq, map[string]string{
        "Last_hop": beginning.Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "0",
    }))
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": beginning.Ip,
        "Seqnum": "1",
        "Metric": "1",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "1",
    }))

    /* check node 2 */
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_received_rreq, map[string]string{
        "Last_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "1",
    }))
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": riot_line[1].Ip,
        "Seqnum": "1",
        "Metric": "2",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_sent_rreq, map[string]string{
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "2",
    }))

    /* check node 3 (aka the end) */
    end.Channels.Expect_JSON(make_JSON_str(template_received_rreq, map[string]string{
        "Last_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Targ_addr": end.Ip,
        "Orig_seqnum": "1",
        "Metric": "2",
    }))
    end.Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": beginning.Ip,
        "Next_hop": riot_line[2].Ip,
        "Seqnum": "1",
        "Metric": "3",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))
    /* And send a RREP back */
    end.Channels.Expect_JSON(make_JSON_str(template_sent_rrep, map[string]string{
        "Next_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
    }))

    /* check node 2 */
    /* TODO FIXME targ_addr_seqnum is *sometimes* 2 because apparently
     * weird RREQs are sent out before the experiment, screwing up the targaddr seqnum
     * and I haven't figured out why yet. So some of these tests may fail, cascading
     * into a whole failure of tests. Yay! :))) */
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_received_rrep, map[string]string{
        "Last_hop": end.Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
        "Targ_seqnum": "2",
    }))
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": end.Ip,
        "Seqnum": "2",
        "Metric": "1",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))
    riot_line[2].Channels.Expect_JSON(make_JSON_str(template_sent_rrep, map[string]string{
        "Next_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
    }))

    /* check node 1 */
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_received_rrep, map[string]string{
        "Last_hop": riot_line[2].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
        "Targ_seqnum": "2",
    }))
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": riot_line[2].Ip,
        "Seqnum": "2",
        "Metric": "2",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_sent_rrep, map[string]string{
        "Next_hop": beginning.Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
    }))

    /* check node 0 (aka the beginning) */
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_received_rrep, map[string]string{
        "Last_hop": riot_line[1].Ip,
        "Orig_addr": beginning.Ip,
        "Orig_seqnum": "1",
        "Targ_addr": end.Ip,
        "Targ_seqnum": "2",
    }))
    riot_line[1].Channels.Expect_JSON(make_JSON_str(template_added_rt_entry, map[string]string{
        "Addr": end.Ip,
        "Next_hop": riot_line[1].Ip,
        "Seqnum": "2",
        "Metric": "3",
        "State": strconv.Itoa(ROUTE_STATE_ACTIVE),
    }))

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