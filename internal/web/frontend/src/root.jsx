import { Outlet } from "react-router-dom";
import React from "react";
import {Box} from "grommet";

const AppBar = (props) => (
    <Box
        tag='header'
        direction='row'
        align='center'
        justify='between'
        background='brand'
        pad={{ left: 'medium', right: 'small', vertical: 'small' }}
        elevation='medium'
        style={{ zIndex: '1' }}
        {...props}
    />
);

function Root() {
    return (
        <>
            <AppBar>HomeThing Manager</AppBar>
            <div id="detail">
                <Outlet />
            </div>
        </>
    );
}
export { Root }