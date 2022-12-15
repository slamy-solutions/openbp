import React, { createRef, useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import Container from '@mui/material/Container';
import TextField, { TextFieldProps } from '@mui/material/TextField';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { LoginWithPasswordRequest, useLoginWithPasswordMutation } from './api.slice'
import { login as setLoginData } from './auth.slice'
import { useDispatch } from 'react-redux';

function Login() {
  const dispatch = useDispatch()

  const errRef = createRef<HTMLInputElement>()

  const [userLogin, setUserLogin] = useState('')
  const [userPassword, setUserPassword] = useState('')
  const [errMsg, setErrMsg] = useState('')

  const [login, { isLoading }] = useLoginWithPasswordMutation()

  useEffect(() => {
    setErrMsg('')
  }, [userLogin, userPassword])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const loginData = {
        login: userLogin,
        password: userPassword
      } as LoginWithPasswordRequest

      const userData = await login(loginData).unwrap()
      dispatch(setLoginData({ ...userData, login: userLogin }))
      setUserLogin('')
      setUserPassword('')
    } catch (err) {
      setErrMsg('Failed to log in.')
      errRef.current?.focus()
    }
  };

  const handleLoginInput = (e: React.ChangeEvent<HTMLInputElement>) => setUserLogin(e.target.value)
  const handlePasswordInput = (e: React.ChangeEvent<HTMLInputElement>) => setUserPassword(e.target.value)

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
            <Typography component="h1" variant="h5">
            User sign in
            </Typography>
            <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
            <TextField
                margin="normal"
                required
                fullWidth
                id="login"
                label="Login"
                name="login"
                autoComplete="login"
                autoFocus
                value={userLogin}
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
                value={userPassword}
                onChange={handlePasswordInput}
            />
            <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                size="large"
            >
                Sign In
            </Button>
            </Box>
        </Box>
    </Container>
  );
}

export default Login;