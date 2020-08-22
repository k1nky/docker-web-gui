import React from 'react'
import { Badge } from 'evergreen-ui'
import { connect } from 'react-redux'


class ContainerStat extends React.PureComponent {
  
  bytesToString(value) {
    var metrics = ["B", "KB", "MB", "GB"]
    var str = value + metrics[0]
    for (var i = 0; i < metrics.length; i++) {
      value = Math.floor(value / 1024)
      if (value == 0) break      
      str = value + metrics[i+1]
    }
    return str
  }

  renderBadges () {
    const { stats, containerID } = this.props
    const data = stats
      .find(n => n.id === containerID)
    return data 
      ? <>
        <Badge backgroundColor="#deebf7" fontWeight="bold" borderRadius={16} paddingLeft={10} fontSize={11} paddingRight={10} marginLeft={10} marginTop={3}>
          cpu {parseFloat(data.cpu_percentage).toFixed(2)} %
        </Badge>
        <Badge backgroundColor="#ebe7f8" fontWeight="bold" borderRadius={16} paddingLeft={10} fontSize={11} paddingRight={10} marginLeft={10} marginTop={3}>
          ram {this.bytesToString(data.memory_usage)} / {this.bytesToString(data.memory_limit)}
        </Badge>
        <Badge backgroundColor="#ebe7f8" fontWeight="bold" borderRadius={16} paddingLeft={10} fontSize={11} paddingRight={10} marginLeft={10} marginTop={3}>
          net {this.bytesToString(data.network_io[0])} / {this.bytesToString(data.network_io[1])}
        </Badge>
      </>
      : null
  }

  render () {
    return this.renderBadges()
  }
}

const mapStateToProps = state => {
  return {
    stats: state.stats.containerStats
  }
}

export default connect(mapStateToProps, null)( ContainerStat )