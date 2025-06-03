import { useState, useEffect } from "react";
import axios from "axios";

function App() {
  const [propsList, setPropsList] = useState([]);
  const [formData, setFormData] = useState({ address: "", city: "", zip: "" });
  const [loading, setLoading] = useState(false);
  const API = import.meta.env.DEV
    ? "http://localhost:8080"
    : "https://your-production-api.example.com";

  // 1. Load all properties
  const loadProperties = async () => {
    try {
      setLoading(true);
      const resp = await axios.get(`${API}/properties`);
      setPropsList(resp.data);
    } catch (err) {
      console.error(err);
      alert("Error loading properties");
    } finally {
      setLoading(false);
    }
  };

  // 2. Handle form submit
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.address || !formData.city || !formData.zip) {
      alert("All fields are required");
      return;
    }
    try {
      await axios.post(`${API}/properties`, formData);
      setFormData({ address: "", city: "", zip: "" });
      loadProperties();
    } catch (err) {
      console.error(err);
      alert("Error creating property");
    }
  };

  // 3. On mount, load properties
  useEffect(() => {
    loadProperties();
  }, []);

  return (
    <div style={{ maxWidth: 600, margin: "auto", padding: 20 }}>
      <h1>Properties</h1>
      <form onSubmit={handleSubmit} style={{ marginBottom: 20 }}>
        <div>
          <label>Address:</label>
          <input
            type="text"
            value={formData.address}
            onChange={(e) => setFormData({ ...formData, address: e.target.value })}
            required
          />
        </div>
        <div>
          <label>City:</label>
          <input
            type="text"
            value={formData.city}
            onChange={(e) => setFormData({ ...formData, city: e.target.value })}
            required
          />
        </div>
        <div>
          <label>ZIP:</label>
          <input
            type="text"
            value={formData.zip}
            onChange={(e) => setFormData({ ...formData, zip: e.target.value })}
            required
          />
        </div>
        <button type="submit">Add Property</button>
      </form>

      {loading ? (
        <p>Loading…</p>
      ) : (
        <table border="1" cellPadding="8" cellSpacing="0" width="100%">
          <thead>
            <tr>
              <th>ID</th>
              <th>Address</th>
              <th>City</th>
              <th>ZIP</th>
              <th>Listed At</th>
            </tr>
          </thead>
          <tbody>
            {propsList.map((p) => (
              <tr key={p.ID}>
                <td>{p.ID}</td>
                <td>{p.Address}</td>
                <td>{p.City}</td>
                <td>{p.ZIP}</td>
                <td>{new Date(p.ListingDate).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default App;
