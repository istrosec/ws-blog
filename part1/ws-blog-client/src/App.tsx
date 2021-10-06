import React, {useEffect, useState} from 'react';
import './App.css';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import {
    Box,
    Collapse,
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography
} from "@mui/material";
import {Terminal} from "./Terminal";
import {Agent} from './Agent';

const serverUrl = "http://localhost:8080"
const agentsUrl = `${serverUrl}/api/agents`
const websocketUrl = `${serverUrl.replace("http", "ws")}/api/client/shell`;

function App() {
    const [agents, setAgents] = useState<Agent[]>([])
    useEffect(() => {
        fetch(
            agentsUrl
        )
            .then(response => response.json())
            .then(data => {
                console.log(JSON.stringify(data))
                setAgents(data as Agent[])
            })
    }, [])
    return (
        <div className="App">
            <TableContainer component={Paper}>
                <Table aria-label="agents-table">
                    <TableHead>
                        <TableRow>
                            <TableCell/>
                            <TableCell>User Name</TableCell>
                            <TableCell align="right">Host Name</TableCell>
                            <TableCell align="right">Local IP</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {agents.map((agent, i) => (
                            <AgentRow key={`${agent.hostName}-${i}`} agent={agent}/>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </div>
    );
}

function AgentRow(props: { agent: Agent }) {
    const {agent} = props;
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow sx={{'& > *': {borderBottom: 'unset'}}}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon/> : <KeyboardArrowDownIcon/>}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    {agent.name}
                </TableCell>
                <TableCell align="right">{agent.hostName}</TableCell>
                <TableCell align="right">{agent.localIp}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{paddingBottom: 0, paddingTop: 0}} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{margin: 1}}>
                            <Typography variant="h6" gutterBottom component="div">
                                Command Prompt
                            </Typography>
                            <Terminal agent={agent} server={websocketUrl}></Terminal>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}

export default App;
