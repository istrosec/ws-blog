package com.istrosec.bloh.ws.server

import com.fasterxml.jackson.databind.SerializationFeature
import io.ktor.application.*
import io.ktor.features.*
import io.ktor.http.cio.websocket.*
import io.ktor.jackson.*
import io.ktor.response.*
import io.ktor.routing.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.websocket.*
import java.time.Duration

fun main() {
    embeddedServer(Netty, port = 8080, host = "0.0.0.0") {
        installJson()
        install(WebSockets)
        routing {
            agent()
            client()
        }
    }.start(wait = true)
}

private fun Routing.agent() {
    webSocket("/agent/shell") { // websocketSession
        passShellMessages(this, agents, clients)
    }
}

private fun Routing.client() {
    webSocket("/client/shell") { // websocketSession
        passShellMessages(this, clients, agents)
    }
    get("/agents") {
        call.respond(agents.list())
    }
}