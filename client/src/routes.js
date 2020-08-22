import React, { useState } from 'react'
import {BrowserRouter, Route, Switch} from 'react-router-dom'
import Navbar from './components/NavBar'
import PrivateRoute from './privateRoute';

import ContainerPage from './pages/container.page'
import ImagePage from './pages/image.page'
import CleanupPage from './pages/cleanup.page'
import LoginPage from './pages/login.page'
import { AuthContext } from './utilities/auth'


const Routes = () => {
    const existingTokens = JSON.parse(localStorage.getItem("tokens"));
    const [authTokens, setAuthTokens] = useState(existingTokens);
    const setTokens = (data) => {
        localStorage.setItem("tokens", JSON.stringify(data));
        setAuthTokens(data);
    }
        return (        
        <React.Fragment>
            <AuthContext.Provider value={{ authTokens, setAuthTokens: setTokens }}>
                <BrowserRouter>
                <Navbar/>
                    <Switch>
                        <PrivateRoute path="/" exact component={ContainerPage}/>
                        <Route path="/login" exact component={LoginPage}/>
                        <PrivateRoute path="/containers" exact component={ContainerPage}/>
                        <PrivateRoute path="/images" component={ImagePage}/>
                        <PrivateRoute path="/cleanup" component={CleanupPage}/>
                    </Switch>
                </BrowserRouter>
            </AuthContext.Provider>
        </React.Fragment>
    )
}

export default Routes