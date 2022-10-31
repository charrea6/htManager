import React from 'react';
import { Grommet } from 'grommet';
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";

import { Root } from './root';
import { Main } from './pages/main';
import { Device } from './pages/device';
import { ErrorPage } from './error-page';


const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        index: true,
        element: <Main />,
      },
      {
        path: 'device/:deviceId',
        element: <Device />
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
