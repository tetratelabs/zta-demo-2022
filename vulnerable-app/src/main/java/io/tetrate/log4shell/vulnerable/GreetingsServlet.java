package io.tetrate.log4shell.vulnerable;

import java.io.IOException;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.nimbusds.jose.JWSObject;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import net.minidev.json.JSONObject;

public class GreetingsServlet extends HttpServlet {

    private static final long serialVersionUID = 1L;

    private static final String NAME_CLAIM = "name";

    private static final String GROUP_CLAIM = "https://zta-demo/group";

    private static Logger log = LogManager.getLogger(GreetingsServlet.class.getName());

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {

        String user = "anonymous";
        String group = null;
        JWSObject token = getToken(request);
        if (token != null) {
            JSONObject payload = token.getPayload().toJSONObject();
            user = (String) payload.get(NAME_CLAIM);
            group = (String) payload.get(GROUP_CLAIM);
            
            // These log lines may trigger the log4shell attach vector if the JWT token contains
            // any malicious claim!
            log.info("token payload: " + payload.toJSONString());
            log.info("user resolved to: " + user);
        }

        response.setContentType("text/plain");
        response.setStatus(HttpServletResponse.SC_OK);
        response.getWriter().println("Welcome, " + user + "!");
        if (group != null) response.getWriter().println("Group: " + group);
        response.getWriter().println("Accessing: " + request.getRequestURI().substring(request.getContextPath().length()));
        if (token != null) response.getWriter().println("\n\nAuthenticated with token:\n" + token.serialize());
    }

    private JWSObject getToken(HttpServletRequest req) {
        String auth = req.getHeader("Authorization");
        if (auth == null) {
            auth = req.getHeader("authorization");
            if (auth == null) {
                log.debug("no authorization header present");
                return null;
            }
        }

        String[] parts = auth.split(" ");
        if (parts.length != 2) {
            log.debug("invalid authorization header value");
            return null;
        }

        try {
            return JWSObject.parse(parts[1]);
        } catch (Exception ex) {
            ex.printStackTrace();
            return null;
        }
    }
}
