import {useParams, useNavigate} from 'react-router-dom';
import {Box, Text, Button, NameValuePair, NameValueList, Heading, Page, PageContent, PageHeader, Anchor} from 'grommet';

export function Device() {
    let { deviceId } = useParams();
    let navigate = useNavigate();
    let toRoot = () => {
        navigate("/");
    }
    return <Page>
        <PageContent>
            <PageHeader title="Device Details" parent={<Anchor label="Back" onClick={toRoot}/>}/>
            <NameValueList>
                <NameValuePair name="id" key="id">{deviceId}</NameValuePair>
                <NameValuePair name=""></NameValuePair>
            </NameValueList>
        </PageContent>
    </Page>;
}