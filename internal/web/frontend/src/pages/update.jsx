import {useParams, useNavigate} from 'react-router-dom';
import {
    Box,
    Button,
    Page,
    PageContent,
    PageHeader,
    Anchor,
    NameValueList, NameValuePair, Select, Text
} from 'grommet';
import {useEffect, useState} from "react";

export function UpdateDevice({devices}) {
    let { deviceId } = useParams();
    let navigate = useNavigate();
    let toRoot = () => {
        navigate("/device/" + deviceId);
    }
    let device = devices.getDeviceInfo(deviceId);
    let description = device == null ? "": device.description;
    let version = device == null ? "": device.version;

    let update = () => {
        const data = new URLSearchParams();
        data.append("command", "update");
        data.append("version", selectedVersion);
        fetch(`/api/devices/${deviceId}/command`, {method: 'post', body: data}).then((response) =>{
            return response.json();
        }).then((response) => {
            if (response.error !== undefined) {
                setStatus(response.error);
            } else {
                setStatus("");
            }
        })
    }

    const [versions, setVersions] = useState([]);
    const [status, setStatus] = useState("");
    const [selectedVersion, setSelectedVersion] = useState('');

    useEffect(() => {
        const loadVersions = () => {
            fetch(`/api/devices/${deviceId}/update/versions`).then((response) =>{
                return response.json();
            }).then((response) =>{
                setVersions(response.versions);
            })
        };
        loadVersions();
    }, [deviceId]);

    return <Page>
        <PageContent height={"large"}>
            <PageHeader title={description} subtitle={"Update"} parent={<Anchor label="Back" onClick={toRoot}/>}/>
            <NameValueList valueProps={{ width: 'large' }}>
                <NameValuePair name="Current Version">{version}</NameValuePair>
                <NameValuePair name="Available Versions">
                    <Box gap={"small"} direction={"row"}>
                    <Select options={versions} onChange={({value, option}) => setSelectedVersion(value)}/>
                    <Button primary label={"Update"} disabled={selectedVersion === ''} onClick={update}/>
                </Box>
                </NameValuePair>
            </NameValueList>
            <Text color={"red"}>{status}</Text>
        </PageContent>
    </Page>;
}