import InputBar from "./InputBar"
import UrlTable  from "./UrlTable";
import './App.css';
import { useEffect, useState } from "react";

const BASE_URL = import.meta.env.VITE_BASE_URL
console.log("BASE_URL");
console.log(BASE_URL);
console.log("BASE_URL");

function App() {
  const [urlDetailRows, setUrlDetailRows] = useState([]);
  useEffect(()=>{
      fetch(`http://${BASE_URL}/admin/route/resolutions/analytics`,
      {method:"GET"}).then(
          response=> response.json()).then(
          data=>{
          setUrlDetailRows(data)
      }
      )
  },[])

  if (urlDetailRows.length>0){ 
    return (
      <div className="app-container"> {/* Wrap the entire app */}
        <InputBar urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows}></InputBar>
        <UrlTable urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows}></UrlTable>
      </div>  
    );

  }else{
    return(
      <div className="app-container"> {/* Wrap the entire app */}
      <InputBar urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows}></InputBar>
      </div>  
    )
  }

}

export default App