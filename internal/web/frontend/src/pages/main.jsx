import React, {useEffect} from 'react';
import { useNavigate } from "react-router-dom";

import {Box, DataTable, Text, Button, PageHeader, PageContent, Page} from 'grommet';
import { Refresh } from 'grommet-icons';

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

const RefreshButton = ({onRefresh}) => {
    const refreshAfter = 30;
    const [timeLeft, setTimeLeft] = React.useState(refreshAfter);

    useEffect(() => {
        // exit early when we reach 0
        if (!timeLeft) {
            onRefresh();
            setTimeLeft(refreshAfter);
            return;
        }

        // save intervalId to clear the interval when the
        // component re-renders
        const intervalId = setInterval(() => {
            setTimeLeft(timeLeft - 1);
        }, 1000);

        // clear interval on re-render to avoid memory leaks
        return () => clearInterval(intervalId);
        // add timeLeft as a dependency to re-rerun the effect
        // when we update it
    }, [timeLeft, onRefresh]);
    const onClick = () => {
        onRefresh();
        setTimeLeft(refreshAfter);
    }
    return <Box direction="row">
        <Text size={"xsmall"} alignSelf={"end"} margin={{right: "small"}}>Refreshing in {timeLeft}s...</Text>
        <Button label="Refresh" onClick={onClick} icon={<Refresh/>}/>
    </Box>
}

export const Main = () => {
      const [sort, setSort] = React.useState({
            property: 'id',
            direction: 'desc',
          });
      const [data, setData] = React.useState([]);
      const navigate = useNavigate();

      const loadData = () => {
          fetch("/api/devices").then((response) =>{
              return response.json();
          }).then((response) =>{
              setData((response));
          })
      };

      useEffect(() => {
          loadData();
      }, []);

      const rowClicked = (event) => {
        navigate('/device/' + event.datum.id)
      }

      return (
          <Page>
              <PageContent>
                  <PageHeader title="Devices" actions={<Box align="end">
                      <RefreshButton onRefresh={loadData}/>
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
