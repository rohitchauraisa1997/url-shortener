import InputBar from "./InputBar"
import UrlTable  from "./UrlTable";
import './App.css';
import { useEffect, useState } from "react";

function App() {
  const [urlDetailRows, setUrlDetailRows] = useState([]);
  useEffect(()=>{
      fetch("http://localhost:3000/admin/route/resolutions/analytics",
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