import { useState } from 'react'

import Container from '@mui/material/Container';
import TextField, { TextFieldProps } from '@mui/material/TextField';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { useInitRootUserMutation } from '../api.slice'

function BootstrapRootUserComponent() {
    const [rootLogin, setRootLogin] = useState('')
    const [rootPassword, setRootPassword] = useState('')

    const [initRootUser, { isLoading }] = useInitRootUserMutation()

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const initResponse = await initRootUser({
            login: rootLogin,
            password: rootPassword
        }).unwrap()
    };
    
    const handleLoginInput = (e: React.ChangeEvent<HTMLInputElement>) => setRootLogin(e.target.value)
    const handlePasswordInput = (e: React.ChangeEvent<HTMLInputElement>) => setRootPassword(e.target.value)

    return (
        <Box component="form" onSubmit={handleSubmit} noValidate sx={{  }} >
            This will create new user with provided login and password. This user will have full access to the system.
            <TextField
                margin="normal"
                required
                fullWidth
                id="login"
                label="Login"
                name="login"
                autoComplete="login"
                autoFocus
                value={rootLogin}
                onChange={handleLoginInput}
            />
            <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
                value={rootPassword}
                onChange={handlePasswordInput}
            />
            <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                size="large"
                disabled={!isLoading}
            >
                Initialize root user
            </Button>
        </Box>
    )
}

export default BootstrapRootUserComponent