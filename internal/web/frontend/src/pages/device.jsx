import {useParams, useNavigate} from 'react-router-dom';
import {
    Box,
    Button,
    Text,
    NameValuePair,
    NameValueList,
    Page,
    PageContent,
    PageHeader,
    Anchor,
    Meter,
    Table,
    TableBody,
    TableHeader,
    TableRow,
    TableCell
} from 'grommet';
import {Update, Upload, Edit, Trash} from "grommet-icons";
import {useEffect, useState} from "react";
import * as dayjs from "dayjs";
import * as relativeTime from "dayjs/plugin/relativeTime";
import * as humanizeDuration from "humanize-duration";

dayjs.extend(relativeTime);

function LastSeen({lastSeen}) {
    const [seen, setSeen] = useState("");
    useEffect(() => {
        if (lastSeen == null) {
            return;
        }
        setSeen(dayjs(lastSeen).fromNow());
        const intervalId = setInterval(() => {
            setSeen(dayjs(lastSeen).fromNow());
        }, 30000);

        return () => clearInterval(intervalId);
    }, [lastSeen]);
    return <Text size={"xsmall"}>{seen}</Text>
}

function MemorySizeText({bytes, label}) {
    let humanizeSize = (b) => {
        if (b > 1024) {
            return Math.floor(bytes / 1024);
        }
        return b;
    }

    let humanizeSizeSuffix = (b) => {
        if (b > 1024) {
            return "KiB";
        }
        return "Bytes";
    }
    return <Box direction="row" align="center"><Text size="large">{humanizeSize(bytes)}</Text><Text size="small">{humanizeSizeSuffix(bytes)} {label}</Text></Box>;
}

function MemoryInfo({free, low}) {
    const [memoryUsage, setMemoryUsage] = useState(false);
    return (
        <Box direction="row">
            <Meter direction="horizontal" max={80*1024} values={[{value: low, highlight: false, onHover: (over) => {
                    setMemoryUsage(over );
                    },}, {value: free - low}]}/>
            <MemorySizeText bytes={memoryUsage ? low : free} label={memoryUsage ? "min free": "free"}/>
        </Box>);
}

function AllTopics({alltopics, values}) {
    let data = Object.entries(alltopics).flatMap(([element, items]) => {
        let result = [];
        let first = true;
        for (const [item, type] of Object.entries(items.pub)) {
            let itemValues = values[element];
            let value = "";
            if (itemValues !== undefined) {
                value = itemValues[item]
                switch (type) {
                    case 6:
                        value = `${value} °c`;
                        break;
                    case 7:
                        value = `${value} %RH`;
                        break;
                    case 8:
                        value = `${value} KPa`;
                        break;
                    default:
                        break;
                }
            }
            result.push({element: first ? element : "", 'item': item, 'value': value });
            first = false;
        }
        return result;
    });
    return <Table>
        <TableHeader>
            <TableRow>
                <TableCell>Element</TableCell>
                <TableCell>Item</TableCell>
                <TableCell>Value</TableCell>
            </TableRow>
        </TableHeader>
        <TableBody>
            { data.map( ({element, item, value}) =>
                    <TableRow>
                        <TableCell scope="row">{element}</TableCell>
                        <TableCell>{item}</TableCell>
                        <TableCell align="right">{value}</TableCell>
                    </TableRow>
            )}
        </TableBody>
    </Table>;
}

export function Device({devices}) {
    let { deviceId } = useParams();
    let navigate = useNavigate();
    let toRoot = () => {
        navigate("/");
    }
    const [info, setInfo] = useState({capabilities:[]});
    const [diag, setDiag] = useState({lastSeen: null, uptime: "", memInfo: {free: 0, low: 0}});
    const [topics, setTopics] = useState({});
    const [values, setValues] = useState({});

    let reboot = () => {
        const data = new URLSearchParams();
        data.append("command", "restart")
        fetch(`/api/devices/${deviceId}/command`, {method: 'post', body: data}).then((response) =>{
            return response.json();
        })
    }

    useEffect(() => {
        devices.selectDevice(deviceId, (msg, data) => {
            switch (msg) {
                case 'info':
                    setInfo(data);
                    break
                case 'diag':
                    setDiag({
                        uptime: humanizeDuration(data.uptime * 1000),
                        lastSeen: data.lastSeen,
                        memInfo: data.mem
                    });
                    break;
                case 'topics':
                    setTopics(data.topics);
                    break;
                case 'values':
                    setValues(data);
                    break;
                case 'value':
                    let newValues = {...values};
                    let itemValues = newValues[data.topic_path[0]];
                    if (itemValues === undefined) {
                        itemValues = newValues[data.topic_path[0]] = {}
                    }
                    itemValues[data.topic_path[1]] = data.value;
                    setValues(newValues);
                    break;
                default:
                    break;
            }
        })
        return () => { devices.unselectDevice(deviceId)};
    }, [deviceId, devices]);

    return <Page>
        <PageContent>
            <PageHeader title={info.description} parent={<Anchor label="Back" onClick={toRoot}/>} actions={<Box direction="row" gap="xsmall">
                <Button plain={false} icon={<Update/>} title={"Reboot"} onClick={reboot}/>
                <Button plain={false} icon={<Upload/>} title={"Update"}/>
                <Button plain={false} icon={<Edit/>} title={"Edit Profile"} onClick={ ()=>{ navigate(`/device/${deviceId}/profile`);} }/>
                <Button plain={false} icon={<Trash/>} title={"Delete device"}/>
            </Box> }/>
            <NameValueList valueProps={{ width: 'large' }}>
                <NameValuePair name="ID">{deviceId}</NameValuePair>
                <NameValuePair name="Version">{info.version}</NameValuePair>
                <NameValuePair name="Capabilities">{info.capabilities.join(', ')}</NameValuePair>
                <NameValuePair name="IP Address"><a href={"http://" + info.ip_addr}>{info.ip_addr}</a></NameValuePair>
                <NameValuePair name="Uptime">{diag.uptime}<LastSeen lastSeen={diag.lastSeen}/></NameValuePair>
                <NameValuePair name="Memory Free"><MemoryInfo free={diag.memInfo.free} low={diag.memInfo.low}/> </NameValuePair>
                <NameValuePair name="Publish Topics"><AllTopics alltopics={topics} values={values}></AllTopics></NameValuePair>
            </NameValueList>
        </PageContent>
    </Page>;
}