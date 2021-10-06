package com.istrosec.blog.ws.server

import io.ktor.http.cio.websocket.*
import kotlinx.coroutines.isActive
import java.util.concurrent.ConcurrentHashMap

data class Agent(
    val name: String,
    val hostName: String,
    val localIp: String
)

class ShellSessions(val name: String) {
    private val storage = ConcurrentHashMap<Agent, WebSocketSession>()

    suspend fun register(agent: Agent, webSocketSession: WebSocketSession) {
        val old = storage.put(agent, webSocketSession)
        old?.close(CloseReason(CloseReason.Codes.GOING_AWAY, "A new connection has been opened."))
    }

    suspend fun unregister(agent: Agent) {
        val old = storage.remove(agent)
        old?.close(CloseReason(CloseReason.Codes.NORMAL, "Shell connection has ended"))
    }

    operator fun get(agent: Agent): WebSocketSession? {
        val ws = storage[agent] ?: return null
        if (!ws.isActive) {
            storage.remove(agent)
            return null
        }
        return ws
    }

    fun list(): List<Agent> {
        return storage.keys.toList()
    }
}

val agents = ShellSessions("AgentWS")

val clients = ShellSessions("ClientWS")