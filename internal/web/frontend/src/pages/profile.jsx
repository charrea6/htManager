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

export function EditProfile() {
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

    const [info, setInfo] = useState({capabilities:[]});
    const [profile, setProfile] = useState("");
    const [status, setStatus] = useState("");
    useEffect(() => {
        const loadInfo = () => {
            fetch(`/api/devices/${deviceId}/info`).then((response) =>{
                return response.json();
            }).then((response) =>{
                setInfo(response);
            })
        };

        const loadProfile = () => {
            fetch(`/api/devices/${deviceId}/profile`).then((response) =>{
                return response.json();
            }).then((response) =>{
                setProfile(response.profile);
            })
        };
        loadInfo();
        loadProfile();
    }, [deviceId]);

    return <Page>
        <PageContent height={"large"}>
            <PageHeader title={info.description} parent={<Anchor label="Back" onClick={toRoot}/>} actions={<Box direction="row" gap="xsmall">
                <Text alignSelf="center">{status}</Text>
                <Button plain={false} icon={<Upload/>} title={"Update"} onClick={updateProfile}/>
            </Box> }/>
            <TextArea value={profile} fill={true} onChange={(event) => { setProfile(event.target.value) }}/>
        </PageContent>
    </Page>;
}