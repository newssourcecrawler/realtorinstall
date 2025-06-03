import { runtime } from '@wailsapp/runtime2'
import { ListProperties, CreateProperty } from './wailsjs/go/services/PropertyService'

console.log("main.js loaded")

document.addEventListener('DOMContentLoaded', () => {
  console.log("DOM fully loaded, attaching event listeners")

  const content = document.getElementById('content')
  const btnList = document.getElementById('btnListProps')
  const btnAdd  = document.getElementById('btnAddProp')

  btnList.addEventListener('click', async () => {
    console.log("List All Properties clicked")
    try {
      // Use the exact exported name: PropertyService (capital P)
      const props = await ListProperties()
      if (props.length === 0) {
        content.innerHTML = '<p>No properties found.</p>'
        return
      }
      let html = '<h2>All Properties</h2>' +
                 '<table><tr><th>ID</th><th>Address</th><th>City</th><th>ZIP</th></tr>'
      props.forEach(p => {
        html += `<tr><td>${p.ID}</td><td>${p.Address}</td><td>${p.City}</td><td>${p.ZIP}</td></tr>`
      })
      html += '</table>'
      content.innerHTML = html
    } catch (err) {
      content.innerHTML = `<p style="color:red;">Error: ${err}</p>`
    }
  })

  btnAdd.addEventListener('click', () => {
    console.log("Add New Property clicked")
    content.innerHTML = `
      <h2>Add New Property</h2>
      <form id="addPropForm">
        <label>Address: <input type="text" id="addr" required /></label><br/><br/>
        <label>City: <input type="text" id="city" required /></label><br/><br/>
        <label>ZIP: <input type="text" id="zip" required /></label><br/><br/>
        <label>Location Code: <input type="text" id="loccode" /></label><br/><br/>
        <label>Size (ft²): <input type="number" id="size" /></label><br/><br/>
        <label>Base Price (USD): <input type="number" id="price" /></label><br/><br/>
        <button type="submit">Save</button>
      </form>
      <div id="result"></div>
    `
    const form = document.getElementById('addPropForm')
    form.addEventListener('submit', async (evt) => {
      evt.preventDefault()
      console.log("Submitting Add New Property form")
      const newProp = {
        Address:      document.getElementById('addr').value,
        City:         document.getElementById('city').value,
        ZIP:          document.getElementById('zip').value,
        LocationCode: document.getElementById('loccode').value,
        SizeSqFt:     parseFloat(document.getElementById('size').value) || 0,
        BasePriceUSD: parseFloat(document.getElementById('price').value) || 0,
      }
      try {
        const id = await CreateProperty(newProp)
        document.getElementById('result').innerHTML =
          `<p style="color:green;">Created property with ID: ${id}</p>`
      } catch (err) {
        document.getElementById('result').innerHTML =
          `<p style="color:red;">Error: ${err}</p>`
      }
    })
  })

  // Optional: let Wails know the frontend is ready
  runtime.EventsOn("frontend:ready", () => {
    console.log("Frontend is ready (Wails runtime)")
  })
})
