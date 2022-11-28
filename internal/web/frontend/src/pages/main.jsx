import React, {useEffect} from 'react';
import { useNavigate } from "react-router-dom";

import {Box, DataTable, Text, PageHeader, PageContent, Page} from 'grommet';

const columns = [
    { header: "", property: "lastSeen", search: false, size: "xsmall", align: "center", render: (data) => {
        return <Alive lastSeen={data.lastSeen}/>
        }},
    { header: "UID", property: "id", search: true, primary: true, render: (data) => {
        return <pre>{data.id}</pre>;
        }},
    { header: "Description", property: "description", search: true},
    { header: "Version", property: "version" , search: true}
];

const Alive = (props) => {
    let lastSeen = Date.parse(props.lastSeen);
    let now = Date.now();
    let initial = "";
    if (now - lastSeen > 45000)
        initial = "💀";
    const [alive, setAlive] = React.useState(initial);

    useEffect(() => {
        const updateAlive = () => {
            let lastSeen = Date.parse(props.lastSeen);
            let now = Date.now();

            if (now - lastSeen > 45000) {
                setAlive("💀");
            } else {
                setAlive(" ");
            }
        }

        updateAlive();
        const timer = setTimeout(updateAlive, 1000);
            return () => {
                clearTimeout(timer);
            }
        }, [props.lastSeen]);
    return <Text>{alive}</Text>;
}

export const Main = ({devices}) => {
      const [sort, setSort] = React.useState({
            property: 'id',
            direction: 'desc',
          });
      const [data, setData] = React.useState([]);
      const navigate = useNavigate();

      useEffect(() => {
          setData(devices.devices);

          devices.deviceListUpdated = (list) => {
              setData(list);
          }
          return () => { devices.deviceListUpdated = null; }
      }, [devices]);

      const rowClicked = (event) => {
        navigate('/device/' + event.datum.id)
      }

      return (
          <Page>
              <PageContent>
                  <PageHeader title="Devices" actions={<Box align="end">
                  </Box>}/>
                <Box fill="horizontal">
                  <DataTable
                    columns={columns}
                    data={data}
                    sort={sort}
                    onSort={setSort}
                    onClickRow={rowClicked}
                  />
                </Box>

              </PageContent>
          </Page>
      );
    };
