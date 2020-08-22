import React, { useState } from "react";
import { Card, Form, Input, Button } from '../components/AuthForm';
import { request } from '../utilities/request';
import { Redirect } from 'react-router-dom';
import { useAuth } from "../utilities/auth";

function LoginPage() {
    const [isLoggedIn, setLoggedIn] = useState(false);
    const [isError, setIsError] = useState(false);
    const [userName, setUserName] = useState("");
    const [password, setPassword] = useState("");
    const { setAuthTokens } = useAuth();

    function postLogin() {
        request('post', 'login', {
            userName,
            password
        }).then(result => {
            if (result.status == 200) {
                setAuthTokens(result.data)
                setLoggedIn(true)
            } else {
                setIsError(true)
            }
        }).catch(e => {
            setIsError(true)
        })
    }
    if (isLoggedIn) {
        return <Redirect to="/"/>;
    }

  return (
    <Card>
      <Form>
        <Input
            type="username"
            placeholder="username"
            value={userName}
            onChange={e => {
                setUserName(e.target.value)
            }}
        />
        <Input
            type="password"
            placeholder="password"
            value={password}
            onChange={e => {
                setPassword(e.target.value)
            }}
        />
        <Button onClick={postLogin}>Login</Button>
      </Form>
    </Card>
  );
}

export default LoginPage;