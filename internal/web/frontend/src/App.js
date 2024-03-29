import React from 'react';
import { Grommet } from 'grommet';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import { Root } from './root';
import { Main } from './pages/main';
import { Device } from './pages/device';
import { EditProfile } from './pages/profile';
import { ErrorPage } from './error-page';
import {DeviceList} from "./Devices";
import {UpdateDevice} from "./pages/update";
import {DeleteDevice} from "./pages/delete";


const deviceList = new DeviceList();

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        index: true,
        element: <Main devices={deviceList}/>,
      },
      {
        path: 'device/:deviceId',
        element: <Device devices={deviceList}/>
      },
      {
        path: 'device/:deviceId/profile',
        element: <EditProfile devices={deviceList}/>
      },
      {
        path: 'device/:deviceId/update',
        element: <UpdateDevice devices={deviceList}/>
      },
      {
        path: 'device/:deviceId/delete',
        element: <DeleteDevice devices={deviceList}/>
      }
    ]
  },

]);

function App() {
  return (
    <Grommet plain>
      <RouterProvider router={router} />
    </Grommet>
  );
}

export default App;
