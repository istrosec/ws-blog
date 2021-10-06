import React, {FC, useEffect, useState} from "react";
import useWebSocket from 'react-use-websocket';
import {IconButton, TextField} from "@mui/material";
import {SendOutlined} from "@mui/icons-material";
import {Agent} from "./Agent";

const maxMessaheHistorySize = 5;

export interface TerminalProps {
    server: string
    agent: Agent
}

export const Terminal: FC<TerminalProps> = (props: TerminalProps) => {
    const {server, agent} = props;
    const [messageHistory, setMessageHistory] = useState<string[]>([]);
    const [inputtedMessage, setInputtedMessage] = useState<string>('');
    const [greetingSent, setGreetingSent] = useState<boolean>(false)
    const {
        sendMessage,
        lastMessage,
    } = useWebSocket(server);

    useEffect(() => {
        lastMessage && setMessageHistory(prev => {
            if (prev.length > maxMessaheHistorySize) {
                let curr = new Array<string>(prev.length)
                for (let i = 1; i < prev.length; i++) {
                    curr[i - 1] = prev[i]
                }
                curr[prev.length - 1] = lastMessage.data
                return curr
            } else {
                return prev.concat(lastMessage.data)
            }
        });
    }, [lastMessage]);

    const onInputChange = (event: any) => {
        setInputtedMessage(event.target.value)
    }

    const handleSendCommand = () => {
        if (!greetingSent) {
            sendMessage(JSON.stringify(agent))
            setGreetingSent(true)
        }
        console.log(`CMD: ${inputtedMessage}`)
        sendMessage(inputtedMessage)
    }
    let currentDir = 'C:\\>';
    if (messageHistory.length > 0) {
        let lastMessageFromHistory = messageHistory[messageHistory.length - 1];

        const workingDirDelimiterPosition = lastMessageFromHistory.indexOf('> ')
        if (workingDirDelimiterPosition > 0) {
            currentDir = `${lastMessageFromHistory.substring(0, workingDirDelimiterPosition)}>`;
        }
    }
    return (
        <div>
            <div className={"shell"}>
                <ul className="commands">
                    {messageHistory.map((message, i) => (
                        message ? <li key={i}>{message}<p/></li> : null
                    ))}
                </ul>
            </div>
            <>
                <TextField id="filled-basic" margin="dense" fullWidth={true} multiline={true} label={currentDir} variant="filled"
                           value={inputtedMessage} onChange={onInputChange}
                           InputProps={{
                               endAdornment: <IconButton color="primary"
                                                         onClick={handleSendCommand}><SendOutlined/></IconButton>
                           }}/>

            </>

        </div>
    );
}