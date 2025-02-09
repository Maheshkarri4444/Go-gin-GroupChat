const backapi = "http://localhost:4000";

const Allapi = {
  signup: {
    url: `${backapi}/api/auth/signup`,
    method: "POST",
  },
  login: {
    url: `${backapi}/api/auth/login`,
    method: "POST",
  },
  logout: {
    url: `${backapi}/api/auth/logout`,
    method: "POST",
  },
  checkauth: {
    url: `${backapi}/api/auth/check-auth`,
    method: "GET",
  },
  messages: {
    url: `${backapi}/api/message/messages`,
    method: "GET",
  },
  sendMessage: {
    url: `${backapi}/api/message/send`,
    method: "POST",
  }
};

export default Allapi