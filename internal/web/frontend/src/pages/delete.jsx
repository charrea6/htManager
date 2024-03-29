import {useParams, useNavigate} from 'react-router-dom';
import {
    Box,
    Button,
    Page,
    PageContent,
    PageHeader,
    Anchor,
    Text
} from 'grommet';
import {useState} from "react";

export function DeleteDevice({devices}) {
    let { deviceId } = useParams();
    let navigate = useNavigate();
    let toDevice = () => {
        navigate("/device/" + deviceId);
    }
    let device = devices.getDeviceInfo(deviceId);
    let description = device == null ? "": device.description;
    const [status, setStatus] = useState("");

    let deleteDevice = () => {
        fetch(`/api/devices/${deviceId}`, {method: 'delete'}).then((response) =>{
            return response.json();
        }).then((response) => {
            if (response.error !== undefined) {
                setStatus(response.error);
            } else {
                navigate("/");
            }
        })
    }

    return <Page>
        <PageContent height={"large"}>
            <PageHeader title={description} parent={<Anchor label="Back" onClick={toDevice}/>}/>
            <Text>Are you sure you want to delete this device?</Text>
            <Box direction={"row"} align={"center"}>
                <Button primary label={"Delete"} onClick={deleteDevice}/>
                <Button label={"Cancel"} onClick={toDevice}/>
            </Box>
            <Text color={"red"}>{status}</Text>
        </PageContent>
    </Page>;
}