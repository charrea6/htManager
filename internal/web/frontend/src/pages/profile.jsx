import {useParams, useNavigate} from 'react-router-dom';
import {
    Box,
    Button,
    Page,
    PageContent,
    PageHeader,
    Anchor,
    TextArea,
    Text
} from 'grommet';
import {Upload,} from "grommet-icons";
import {useEffect, useState} from "react";

export function EditProfile({devices}) {
    let { deviceId } = useParams();
    let navigate = useNavigate();
    let toRoot = () => {
        navigate("/device/" + deviceId);
    }

    let updateProfile = () => {
        fetch(`/api/devices/${deviceId}/profile`, {method: 'post', body: profile}).then((response) =>{
            return response.json();
        }).then((response) => {
            if (response.error === undefined) {
                setStatus("Profile updated.");
                navigate("/device/" + deviceId);
            } else {
                setStatus("Profile update failed! " + response.error);
            }
            } ).catch(() => { setStatus('Profile update failed')});
    }

    let device = devices.getDeviceInfo(deviceId);
    let description = device == null ? "": device.description;

    const [profile, setProfile] = useState("");
    const [status, setStatus] = useState("");
    useEffect(() => {
        const loadProfile = () => {
            fetch(`/api/devices/${deviceId}/profile`).then((response) =>{
                return response.json();
            }).then((response) =>{
                setProfile(response.profile);
            })
        };
        loadProfile();
    }, [deviceId]);

    return <Page>
        <PageContent height={"large"}>
            <PageHeader title={description} subtitle={"Edit Profile"} parent={<Anchor label="Back" onClick={toRoot}/>} actions={<Box direction="row" gap="xsmall">
                <Text alignSelf="center">{status}</Text>
                <Button plain={false} icon={<Upload/>} title={"Update"} onClick={updateProfile}/>
            </Box> }/>
            <TextArea value={profile} fill={true} onChange={(event) => { setProfile(event.target.value) }}/>
        </PageContent>
    </Page>;
}