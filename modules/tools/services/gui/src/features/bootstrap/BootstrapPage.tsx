import Container from '@mui/material/Container';
import TextField, { TextFieldProps } from '@mui/material/TextField';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import Paper from '@mui/material/Paper';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import BootstrapRootUserComponent from './components/rootUser'

import { store } from '../../store'
import { useLazyGetStatusQuery, useGetStatusQuery, useInitRootUserMutation } from './api.slice'
import { setState } from './bootstrap.slice'
import { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';

function BootstrapPage() {   
    const dispatch = useDispatch() 
    const navigate = useNavigate()
    const {data: status, error: statusLoadError, isLoading: isStatusLoading} = useGetStatusQuery({})
    //const [initRootUser, { isLoading }] = useInitRootUserMutation()

    useEffect(()=>{
        if (status) {
            dispatch(setState({
                fullyBootstrapped: status.fullyBootstrapped,
                rootUserInitialized: true
            }))
            navigate('/')
        }
    }, [status])

    return (
        <Container component="main" maxWidth="sm">
            <Box
                sx={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: "center",
                alignItems: 'center',
                minHeight: '100vh'
                }}
            >
                <img src='/logo.svg' alt="logo" style={{width: "60%", marginBottom: "42px"}}/>
                <Typography component="h1" variant="h5" sx={{ mb: 5 }}>
                Bootstrapping system
                </Typography>

                {isStatusLoading && ( <CircularProgress /> )}

                {!isStatusLoading && !store.getState()['features/bootstrap'].rootUserInitialized && (
                    <Paper sx={{ p: 2, justifyContent: "center", alignItems: 'center', flexDirection: 'column', display: 'flex' }}>
                        <Typography component="h2" variant="h6">
                            Initialize root user
                        </Typography>
                        <BootstrapRootUserComponent />
                    </Paper>
                )}
            </Box>
        </Container>
    )
}

export default BootstrapPage