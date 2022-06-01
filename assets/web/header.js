const headerStyles = `
#head {
    border-bottom: 1px solid #fff;
    background-color: #181818;
}
#head h1 {
    letter-spacing: -4px;
    color: #fff;
    font-size: 32px;
    font-family: 'Times New Roman', Times, serif;
    transform: scale(.5, 1);
    letter-spacing: 1px;
    padding: 0;
    margin: 0;
    display: inline-block;
    position: relative;
    left: -29px;
    margin-right: -55px;
}
#head h1 a {
    color: #fff;
    text-decoration: none;
}
#head h1 span {
    color: #ff3b3b;
}
#head h2 {
    font-style: italic;
    font-weight: normal;
    font-size: 19px;
    display: inline-block;
    margin: 0;
    padding: 0;
    position: relative;
    top: -3px;
}
`;

function fftInjectHeader() {
    let element = document.getElementById('head');
    let pageTitle = element.getAttribute('data-page-title');
    let logoEle = document.createElement('h1');
    let logSpanEle = document.createElement('span');
    logSpanEle.innerText = 'FF';
    logoEle.append(logSpanEle);
    logoEle.append(document.createTextNode('TOOLS'));    
    if (pageTitle) {
        let pageTitleEle = document.createElement('h2');
        pageTitleEle.innerText = pageTitle;
        element.prepend(pageTitleEle);
    }
    element.prepend(logoEle);
    let styleEle = document.createElement('style');
    styleEle.innerHTML = headerStyles;
    element.prepend(styleEle);
    window.removeEventListener('mousemove', fftInjectHeader);
}
window.addEventListener('mousemove', fftInjectHeader);

