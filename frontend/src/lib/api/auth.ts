const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

export async function registerUser(data: {
  email: string;
  fullname: string;
  birthday: string;
  gender: string;
  password: string;
}) {
  const res = await fetch(`${API_BASE}/auth/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  //console.log("Gửi đăng ký:", data);

  const json = await res.json();
  if (!res.ok) throw new Error(json.error || "Registration failed");
  return json;
}

export async function loginUser(data: { email: string; password: string }) {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  const json = await res.json();
  if (!res.ok) throw new Error(json.message || "Login failed");
  return json.access_token;
}
