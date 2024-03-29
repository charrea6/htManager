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

function MemoryInfo({free, low, total}) {
    const [memoryUsage, setMemoryUsage] = useState(false);
    return (
        <Box direction="row">
            <Meter direction="horizontal" max={total} values={[{value: low, highlight: false, onHover: (over) => {
                    setMemoryUsage(over );
                    },}, {value: free - low}]}/>
            <MemorySizeText bytes={memoryUsage ? low : free} label={memoryUsage ? "min free": "free"}/>
        </Box>);
}

function AllTopics({alltopics, values}) {
    let data = Object.entries(alltopics).flatMap(([element, items]) => {
        let result = [];
        let idx = 0;
        for (const [item, type] of Object.entries(items.pub)) {
            let itemValues = values[element];
            let value = "";
            if (itemValues !== undefined) {
                value = itemValues[item]
            }
            let typeStr = "";
            const types = ["Bool", "Int", "Float", "String", "Binary", "Hundredths", "°c", "%RH", "KPa"];
            if (type >= 0 && type < types.length) {
                typeStr = types[type];
            }
            result.push({element: idx === 0 ? element : "", key: `${element}-${idx}`, 'item': item, 'value': value, 'type': typeStr });
            idx ++;
        }
        return result;
    });
    return <Table>
        <TableHeader>
            <TableRow>
                <TableCell>Element</TableCell>
                <TableCell>Item</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Value</TableCell>
            </TableRow>
        </TableHeader>
        <TableBody>
            { data.map( ({element, key, item, value, type}) =>
                    <TableRow key={key}>
                        <TableCell scope="row" >{element}</TableCell>
                        <TableCell>{item}</TableCell>
                        <TableCell>{type}</TableCell>
                        <TableCell>{value}</TableCell>
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
    const [status, setStatus] = useState("");

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
                    setValues((state) => {
                            let newValues = {...state};
                            let itemValues = newValues[data.topic_path[0]];
                            if (itemValues === undefined) {
                                itemValues = newValues[data.topic_path[0]] = {}
                            }
                            itemValues[data.topic_path[1]] = data.value;
                            return newValues
                        }
                    );
                    break;
                case 'status':
                    setStatus(data);
                    break;
                default:
                    break;
            }
        })
        return () => { devices.unselectDevice(deviceId)};
    }, [deviceId, devices]);

    return <Page>
        <PageContent>
            <PageHeader title={info.description} subtitle={"Device Details"} parent={<Anchor label="Back" onClick={toRoot}/>} actions={<Box direction="row" gap="xsmall">
                <Button plain={false} icon={<Update/>} title={"Reboot"} onClick={reboot}/>
                <Button plain={false} icon={<Upload/>} title={"Update"} onClick={ ()=>{ navigate(`/device/${deviceId}/update`);}}/>
                <Button plain={false} icon={<Edit/>} title={"Edit Profile"} onClick={ ()=>{ navigate(`/device/${deviceId}/profile`);} }/>
                <Button plain={false} icon={<Trash/>} title={"Delete"} onClick={ ()=>{ navigate(`/device/${deviceId}/delete`);} }/>
            </Box> }/>
            <NameValueList valueProps={{ width: 'large' }}>
                <NameValuePair name="ID">{deviceId}</NameValuePair>
                <NameValuePair name="Device Type">{info.deviceType}</NameValuePair>
                <NameValuePair name="Version">{info.version}</NameValuePair>
                <NameValuePair name="Capabilities">{info.capabilities.join(', ')}</NameValuePair>
                <NameValuePair name="IP Address"><a href={"http://" + info.ip_addr}>{info.ip_addr}</a></NameValuePair>
                <NameValuePair name="Uptime">{diag.uptime}<LastSeen lastSeen={diag.lastSeen}/></NameValuePair>
                <NameValuePair name="Memory Free"><MemoryInfo free={diag.memInfo.free} low={diag.memInfo.low} total={info.memory}/> </NameValuePair>
                <NameValuePair name="Status">{status}</NameValuePair>
                <NameValuePair name="Publish Topics"><AllTopics alltopics={topics} values={values}></AllTopics></NameValuePair>
            </NameValueList>
        </PageContent>
    </Page>;
}